package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/aquasecurity/libbpfgo"
	"github.com/aquasecurity/libbpfgo/helpers"
)

func main() {

	bpfProgramName := os.Args[1]
	pathToTracingProgram := os.Args[2] //todo: validate path
	symbolName := os.Args[3]

	offset, err := helpers.SymbolToOffset(pathToTracingProgram, symbolName)
	if err != nil {
		log.Fatal(err)
	}

	module, err := libbpfgo.NewModuleFromFile("./dist/weaver.bpf.o") //todo: embed
	if err != nil {
		log.Fatal(err)
	}

	prog, err := module.GetProgram(bpfProgramName)
	if err != nil {
		log.Fatal(err)
	}

	err = module.BPFLoadObject()
	if err != nil {
		log.Fatal(err)
	}

	_, err = prog.AttachUprobe(-1, pathToTracingProgram, uint32(offset))
	if err != nil {
		log.Fatal(err)
	}

	if bpfProgramName == "repeat" {
		repeat()
		os.Exit(0)
	}

	eventsChannel := make(chan []byte)
	rb, err := module.InitRingBuf("events", eventsChannel)
	if err != nil {
		log.Fatal(err)
	}

	rb.Start()

	for {
		fmt.Println("waiting...")
		b := <-eventsChannel
		x := trace_context{}

		stack := b[:50]
		x.Stack = *(*Mem_stack)(unsafe.Pointer(&stack))
		registers := b[56:]

		//TODO: define an interface for 64/32, ARM/AMD
		x.Registers.R15 = int64(binary.LittleEndian.Uint64(registers[0:8]))
		x.Registers.R14 = int64(binary.LittleEndian.Uint64(registers[8:16]))
		x.Registers.R13 = int64(binary.LittleEndian.Uint64(registers[16:24]))
		x.Registers.R12 = int64(binary.LittleEndian.Uint64(registers[24:32]))
		x.Registers.Bp = int64(binary.LittleEndian.Uint64(registers[32:40]))
		x.Registers.Bx = int64(binary.LittleEndian.Uint64(registers[40:48]))
		x.Registers.R11 = int64(binary.LittleEndian.Uint64(registers[48:56]))
		x.Registers.R10 = int64(binary.LittleEndian.Uint64(registers[56:64]))
		x.Registers.R9 = int64(binary.LittleEndian.Uint64(registers[64:72]))
		x.Registers.R8 = int64(binary.LittleEndian.Uint64(registers[72:80]))
		x.Registers.Ax = int64(binary.LittleEndian.Uint64(registers[80:88]))
		x.Registers.Cx = int64(binary.LittleEndian.Uint64(registers[88:96]))
		x.Registers.Dx = int64(binary.LittleEndian.Uint64(registers[96:104]))
		x.Registers.Si = int64(binary.LittleEndian.Uint64(registers[104:112]))
		x.Registers.Di = int64(binary.LittleEndian.Uint64(registers[112:120]))
		x.Registers.Orig_ax = int64(binary.LittleEndian.Uint64(registers[120:128]))
		x.Registers.Ip = int64(binary.LittleEndian.Uint64(registers[128:136]))
		x.Registers.Cs = int64(binary.LittleEndian.Uint64(registers[136:144]))
		x.Registers.Flags = int64(binary.LittleEndian.Uint64(registers[144:152]))
		x.Registers.Sp = int64(binary.LittleEndian.Uint64(registers[152:160]))
		x.Registers.Ss = int64(binary.LittleEndian.Uint64(registers[160:168]))

		jx, err := json.MarshalIndent(x, "", " ")
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println(string(jx))

	}

	rb.Stop()
	rb.Close()
}

func repeat() {
	block := make(chan []byte)
	<-block
}
