package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	bpf "github.com/iovisor/gobpf/bcc"
	"log"
	"math"
	"os"
	"os/signal"
	"text/template"
)

const bpfProgramTextTemplate = `
	#include <uapi/linux/ptrace.h>
	#include <linux/string.h>

	BPF_PERF_OUTPUT(events);

	inline int print_symbol_arg(struct pt_regs *ctx) {
		void* stackAddr = (void*)ctx->sp;
		{{range $arg_index, $arg_element := .Arguments}}
		{{if eq $arg_element.CType "char *" }}
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

func bpfText(context *traceContext) string {
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
func loadUprobeAndBPFModule(context *traceContext) error {

	// Load eBPF filter and uprobe
	filterText := bpfText(context)
	bpfModule := bpf.NewModule(filterText, []string{})
	defer bpfModule.Close()

	uprobeFd, err := bpfModule.LoadUprobe("print_symbol_arg")
	if err != nil {
		return err
	}

	err = bpfModule.AttachUprobe(context.binaryName, context.functionName, uprobeFd, -1)
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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	numberOfArgs := len(context.Arguments)
	var index int
	var dataTypeOfValue goType

	go func() {

		for {
			value := <-channel

			// based on order of value coming in determine what type it is for interpretation
			dataTypeOfValue = context.Arguments[index].goType

			valueString := interpretDataByType(value, dataTypeOfValue)

			outputValue := output{
				Type:  goTypeToString[dataTypeOfValue],
				Value: valueString,
			}

			printOutput(outputValue)

			index++
			index = index % numberOfArgs
		}
	}()

	perfMap.Start()
	<-sig
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
		return fmt.Sprintf("%f", val)
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
		return fmt.Sprintf("%s", data)
	//TODO:
	case STRUCT:
		return "struct interpretation is not yet implemented"
	case POINTER:
		return "pointer interpretation is not yet implemented"
	}

	return ""
}
