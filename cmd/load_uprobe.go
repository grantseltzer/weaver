package main

import (
	bpf "github.com/iovisor/gobpf/bcc"
	"fmt"
	"bytes"
	"log"
	"os"
	"os/signal"
	"text/template"
	"math"
	"encoding/binary"
)

import "C"

const textTemplate = `
	#include <uapi/linux/ptrace.h>
	#include <linux/string.h>

	BPF_PERF_OUTPUT(events);

	inline int print_symbol_arg(struct pt_regs *ctx) {
	
		void* stackAddr = (void*)ctx->sp;
		{{range $arg_index, $arg_element := .Arguments}}
		{{$arg_element.CType}} {{$arg_element.VariableName}};
		void* stackPtr_{{$arg_element.VariableName}} = &stackAddr+{{$arg_element.StartingOffset}};
		{{$arg_element.CType}}* typeStackPtr_{{$arg_element.VariableName}} = ({{$arg_element.CType}}*)stackPtr_{{$arg_element.VariableName}};
		bpf_probe_read(&{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}), stackAddr+{{$arg_element.StartingOffset}}); 
		events.perf_submit(ctx, &{{$arg_element.VariableName}}, sizeof({{$arg_element.VariableName}}));
		{{end}}
		return 0;
	}
`

func bpfText(context *traceContext) string { 
	t := template.New("bpf_text")
	t, err := t.Parse(textTemplate)
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

func createBPFModule(context *traceContext) error {

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
			fmt.Println(valueString) 

			index++
			index = index%numberOfArgs
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

	switch gt{

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
	//TODO:
	case STRING:
		return "string interpretation is not yet implemented"
	case STRUCT:
		return "struct interpretation is not yet implemented"
	case POINTER:
		return "pointer interpretation is not yet implemented"
	}

	return ""
}