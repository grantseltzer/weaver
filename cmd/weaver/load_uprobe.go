package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"text/template"

	bpf "github.com/iovisor/gobpf/bcc"
)

const bpfWithArgsProgramTextTemplate = `
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

		{{if gt .Filters.Pid 0}}
		// apply pid filters
		if (procInfo.pid != {{ .Filters.Pid }}) {
			return 0;
		}
		{{end}}

		// submit process info
		events.perf_submit(ctx, &procInfo, sizeof(procInfo));

		{{if eq .HasArguments true}}

			void* stackAddr = (void*)ctx->sp;

			// [TEMPLATE] Traverse over each argument in this trace context
			{{range $arg_index, $arg_element := .Arguments}}

			    // [TEMPLATE] If this argument is a slice
				{{if eq $arg_element.IsSlice true}}

					// read in bytes for:
					// (8) - array address
					// (8) - array length
					// (8) - array cap
					unsigned long {{$arg_element.VariableName}}_addr;
					unsigned long {{$arg_element.VariableName}}_length;
					unsigned long {{$arg_element.VariableName}}_cap;

					unsigned int i_{{$arg_element.VariableName}};

					for (i_{{$arg_element.VariableName}} = 0; i_{{$arg_element.VariableName}} < {{$arg_element.VariableName}}_length; i_{{$arg_element.VariableName}}++) {
						
					}

				
				// [TEMPLATE] If it's not a slice, but this argument is an array
				{{else if gt $arg_element.ArrayLength 0}}

					unsigned int i_{{$arg_element.VariableName}};
					void* loopAddr_{{$arg_element.VariableName}} = stackAddr+{{$arg_element.StartingOffset}};
					for (i_{{$arg_element.VariableName}} = 0; i_{{$arg_element.VariableName}} < {{$arg_element.ArrayLength}}; i_{{$arg_element.VariableName}}++) {
						
						// [TEMPLATE] This is NOT an array of strings
						{{if ne $arg_element.CType "char *" }}

							{{$arg_element.CType}} {{$arg_element.VariableName}};
							bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), loopAddr_{{$arg_element.VariableName}}); 
							events.perf_submit(ctx, &{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}));
							loopAddr_{{$arg_element.VariableName}} += {{$arg_element.TypeSize}};
						
						// [TEMPLATE] This is an array of strings
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

				// [TEMPLATE] If it's not array, but it's a string
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

				// [TEMPLATE] Any other type besides an array or string
				{{- else }}

					{{$arg_element.CType}} {{$arg_element.VariableName}};
					bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), stackAddr+{{$arg_element.StartingOffset}}); 
					events.perf_submit(ctx, &{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}));
				
				{{- end}}
			
			{{end}}
		{{end}}
		return 0;
	}
`

// bpfText compiles the traceContext into a eBPF program using the above text tempate
func bpfText(context *functionTraceContext) string {
	t := template.New("bpf_text")
	t, err := t.Parse(bpfWithArgsProgramTextTemplate)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	t.Execute(buf, context)

	// Print eBPF text
	debugeBPFLog("%s\n", buf.String())

	return buf.String()
}

// loadUprobeAndBPFModule will, based on the passed traceContext, install the bpf program, attach a uprobe to the specified function
// It then prints results to the designated output stream. Will handle with or without arguments depending on value of 'globalMode'
// This blocks until runtimeContext.Done() triggers
func loadUprobeAndBPFModule(traceContext *functionTraceContext, runtimeContext context.Context, wg *sync.WaitGroup) error {

	defer runtimeContext.Err()

	// Generate eBPF code via text template and load it into a new module
	filterText := bpfText(traceContext)
	bpfModule := bpf.NewModule(filterText, []string{})
	defer bpfModule.Close()

	// Attach the loaded eBPF code to a uprobe'd function specified by the traceContext.FunctionName
	debugLog("Attaching uprobe to %s\n", traceContext.FunctionName)
	uprobeFd, err := bpfModule.LoadUprobe("print_symbol_arg") // name of eBPF function
	if err != nil {
		return err
	}
	err = bpfModule.AttachUprobe(traceContext.binaryName, traceContext.FunctionName, uprobeFd, -1)
	if err != nil {
		return fmt.Errorf("could not attach uprobe to symbol: %s: %s", "test_function", err.Error())
	}

	// Set up bpf perf map to use for output from eBPF to weaver
	table := bpf.NewTable(bpfModule.TableId("events"), bpfModule)
	channel := make(chan []byte)
	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		return err
	}

	if globalMode == PACKAGE_MODE {
		go withoutArgumentsListen(traceContext.FunctionName, channel)
	} else {
		go withArgumentsListen(traceContext, channel)
	}

	wg.Done()
	perfMap.Start()
	<-runtimeContext.Done()
	perfMap.Stop()

	return nil
}

// withArgumentsListen will listen for output from the channel which received output from the eBPF program.
// It reads in process information, followed by associated arguments and prints them
func withArgumentsListen(traceContext *functionTraceContext, rawBytes chan []byte) {

	var (
		output          = output{FunctionName: traceContext.FunctionName}
		numberOfArgs    = len(traceContext.Arguments)
		index           int
		dataTypeOfValue goType
		argOutput       = make([]outputArg, numberOfArgs)
		valueString     string
		outputValue     outputArg
	)

	for {
		value := <-rawBytes
		procInfo := procInfo{}
		err := procInfo.unmarshalBinary(value)
		if err == nil {
			output.ProcInfo = procInfo
			// if err == nil value was proc info struct, else do
			// not fetch next value
			value = <-rawBytes
		}

		// Determine what type it is for interpretation based on order of value coming in
		dataTypeOfValue = traceContext.Arguments[index].goType

		// If this argument is an array
		if traceContext.Arguments[index].ArrayLength > 0 {

			arrayValueString := interpretDataByType(value, dataTypeOfValue)

			for i := 0; i < traceContext.Arguments[index].ArrayLength-1; i++ {
				value := <-rawBytes
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
}

// withoutArgumentsListen will listen for output from the channel which received output from the eBPF program.
// It reads in process information, puts the function name in the output, and prints it
func withoutArgumentsListen(functionName string, rawBytes chan []byte) {
	for {
		value := <-rawBytes
		procInfo := procInfo{}
		err := procInfo.unmarshalBinary(value)
		if err != nil {
			debugLog("could not read in proccess information: %s\n", err.Error())
		}

		output := output{FunctionName: functionName}
		output.ProcInfo = procInfo
		printOutput(output)
	}
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
