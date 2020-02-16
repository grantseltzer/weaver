package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"
	"text/template"

	bpf "github.com/iovisor/gobpf/bcc"
)

const bpfProgramTextTemplate = `
	#include <uapi/linux/ptrace.h>
	#include <linux/string.h>
	#include <linux/fs.h>
	#include <linux/sched.h>

	BPF_PERF_OUTPUT(events);

	struct proc_info_t {
    	u32 pid;  // PID as in the userspace term (i.e. task->tgid in kernel)
    	u32 ppid; // Parent PID as in the userspace term (i.e task->real_parent->tgid in kernel)
    	char comm[TASK_COMM_LEN]; // 16 bytes
	};

	inline int print_symbol_arg(struct pt_regs *ctx) {
		
		// get process info
    	struct task_struct *task;
		struct proc_info_t procInfo = {};
    	task = (struct task_struct *)bpf_get_current_task();
		procInfo.pid = bpf_get_current_pid_tgid() >> 32;
		procInfo.ppid = task->real_parent->tgid;
		bpf_get_current_comm(&procInfo.comm, sizeof(procInfo.comm));

		events.perf_submit(ctx, &procInfo, sizeof(procInfo));

		void* stackAddr = (void*)ctx->sp;
		{{range $arg_index, $arg_element := .Arguments}}

		{{if gt $arg_element.ArrayLength 0}}

		unsigned int i_{{$arg_element.VariableName}};
		void* loopAddr_{{$arg_element.VariableName}} = stackAddr+{{$arg_element.StartingOffset}};
		for (i_{{$arg_element.VariableName}} = 0; i_{{$arg_element.VariableName}} < {{$arg_element.ArrayLength}}; i_{{$arg_element.VariableName}}++) {
			{{if ne $arg_element.CType "char *" }} 
			{{$arg_element.CType}} {{$arg_element.VariableName}};
			bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), loopAddr_{{$arg_element.VariableName}}); 
			events.perf_submit(ctx, &{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}));
			loopAddr_{{$arg_element.VariableName}} += {{$arg_element.TypeSize}};
			{{else}}
			unsigned long {{$arg_element.VariableName}}_length;
			bpf_probe_read(&{{$arg_element.VariableName}}_length, sizeof({{$arg_element.VariableName}}_length), loopAddr_{{$arg_element.VariableName}}+8);
			if ({{$arg_element.VariableName}}_length > 16 ) {
				{{$arg_element.VariableName}}_length = 16;
			}
			unsigned int str_length = (unsigned int){{$arg_element.VariableName}}_length;
			
			// use long double to have up to a 16 character string by reading in the raw bytes
			long double* {{$arg_element.VariableName}}_ptr;
			long double  {{$arg_element.VariableName}};
			bpf_probe_read(&{{$arg_element.VariableName}}_ptr, sizeof({{$arg_element.VariableName}}_ptr), loopAddr_{{$arg_element.VariableName}});
			bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), {{$arg_element.VariableName}}_ptr);
		
			events.perf_submit(ctx, &{{$arg_element.VariableName}}, str_length);
			loopAddr_{{$arg_element.VariableName}} += 16;
			{{end}}
		}
		{{else if eq $arg_element.CType "char *" }}
		unsigned long {{$arg_element.VariableName}}_length;
		bpf_probe_read(&{{$arg_element.VariableName}}_length, sizeof({{$arg_element.VariableName}}_length), stackAddr+{{$arg_element.StartingOffset}}+8);
		if ({{$arg_element.VariableName}}_length > 16 ) {
			{{$arg_element.VariableName}}_length = 16;
		}
		unsigned int str_length = (unsigned int){{$arg_element.VariableName}}_length;
		
		// use long double to have up to a 16 character string by reading in the raw bytes
		long double* {{$arg_element.VariableName}}_ptr;
		long double  {{$arg_element.VariableName}};
		bpf_probe_read(&{{$arg_element.VariableName}}_ptr, sizeof({{$arg_element.VariableName}}_ptr), stackAddr+{{$arg_element.StartingOffset}});
		bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), {{$arg_element.VariableName}}_ptr);
		events.perf_submit(ctx, &{{$arg_element.VariableName}}, str_length);
		{{- else }}
		{{$arg_element.CType}} {{$arg_element.VariableName}};
		bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), stackAddr+{{$arg_element.StartingOffset}}); 
		events.perf_submit(ctx, &{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}));
		{{- end}}
		{{end}}
		return 0;
	}
`

func bpfText(context *functionTraceContext) string {
	t := template.New("bpf_text")
	t, err := t.Parse(bpfProgramTextTemplate)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	t.Execute(buf, context)

	if globalDebug {
		// Print eBPF text
		fmt.Println(buf.String())
	}

	return buf.String()
}

// loadUprobeAndBPFModule will, based on the passed context, install the bpf program and attach a uprobe to the specified function
// It then prints results to the designated output stream.
// This blocks until Ctrl-C or error occurs.
func loadUprobeAndBPFModule(traceContext *functionTraceContext, runtimeContext context.Context) error {

	// Load eBPF filter and uprobe
	filterText := bpfText(traceContext)
	bpfModule := bpf.NewModule(filterText, []string{})
	defer bpfModule.Close()

	uprobeFd, err := bpfModule.LoadUprobe("print_symbol_arg")
	if err != nil {
		return err
	}

	err = bpfModule.AttachUprobe(traceContext.binaryName, traceContext.FunctionName, uprobeFd, -1)
	if err != nil {
		return fmt.Errorf("could not attach uprobe to symbol: %s: %s", "test_function", err.Error())
	}

	// Set up bpf perf map to use for output
	table := bpf.NewTable(bpfModule.TableId("events"), bpfModule)
	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		return err
	}

	numberOfArgs := len(traceContext.Arguments)
	var index int
	var dataTypeOfValue goType
	output := output{FunctionName: traceContext.FunctionName}
	var argOutput = make([]outputArg, numberOfArgs)
	go func() {
		// First sent values are process info
		value := <-channel
		output.ProcInfo = procInfo{}
		err := output.ProcInfo.UnmarshalBinary(value)
		if err != nil {
			fmt.Println("failed to unmarshall process info")
		}

		var valueString string
		var outputValue outputArg
		for {
			value := <-channel

			// Determine what type it is for interpretation based on order of value coming in
			dataTypeOfValue = traceContext.Arguments[index].goType

			// If this argument is an array
			if traceContext.Arguments[index].ArrayLength > 0 {

				arrayValueString := interpretDataByType(value, dataTypeOfValue)

				for i := 0; i < traceContext.Arguments[index].ArrayLength-1; i++ {
					value := <-channel
					valueString = interpretDataByType(value, dataTypeOfValue)
					arrayValueString = arrayValueString + ", " + valueString
				}
				outputValue = outputArg{
					Type:  goTypeToString[dataTypeOfValue] + "_ARRAY",
					Value: arrayValueString,
				}

			} else {
				// This argument is not an array

				valueString = interpretDataByType(value, dataTypeOfValue)

				outputValue = outputArg{
					Type:  goTypeToString[dataTypeOfValue],
					Value: valueString,
				}
			}

			argOutput[index] = outputValue
			index++
			index = index % numberOfArgs

			if index == 0 {
				output.Args = argOutput
				printOutput(output)
			}

		}
	}()

	perfMap.Start()
	<-runtimeContext.Done()
	perfMap.Stop()

	return nil
}

// interpretDataByType takes raw bytes of a value, and returns a string
// where the value is displayed as a type specified by the goType
func interpretDataByType(data []byte, gt goType) string {

	switch gt {

	case INT8, INT16, INT32, UINT8, UINT16, UINT32:
		x1 := binary.LittleEndian.Uint32(data)
		return fmt.Sprintf("%d", int(x1))
	case INT, INT64, UINT, UINT64:
		x1 := binary.LittleEndian.Uint64(data)
		return fmt.Sprintf("%d", int(x1))
	case FLOAT32:
		x1 := binary.LittleEndian.Uint32(data)
		val := math.Float32frombits(x1)
		return fmt.Sprintf(stringfFormat(gt), val)
	case FLOAT64:
		x1 := binary.LittleEndian.Uint64(data)
		val := math.Float64frombits(x1)
		return fmt.Sprintf(stringfFormat(gt), val)
	case BOOL:
		x1 := binary.LittleEndian.Uint32(data)
		if x1 == 0 {
			return "false"
		}
		return "true"
	case BYTE:
		x1 := binary.LittleEndian.Uint32(data)
		return fmt.Sprintf("dec=%d\tchar='%c'", x1, x1)
	case STRING:
		stringValue := strings.SplitN(string(data), "\u0000", 2)
		return fmt.Sprintf("'%s'", stringValue[0])
	case RUNE:
		x1 := binary.LittleEndian.Uint32(data)
		return fmt.Sprintf(stringfFormat(gt), int(x1))
	//TODO:
	case STRUCT:
		return "struct interpretation is not yet implemented"
	case POINTER:
		return "pointer interpretation is not yet implemented"
	}

	return ""
}
