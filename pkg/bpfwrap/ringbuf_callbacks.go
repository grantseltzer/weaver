package bpfwrap

import (
	"C"
	"fmt"
	"unsafe"
)

//export ringbufCallback
func ringbufCallback(ctx unsafe.Pointer, cpu C.int, data unsafe.Pointer, size C.int) {
	fmt.Println(C.GoBytes(data, size))
	//TODO: unmarshalBinary data as littleEndian etc...
}

//export ringbufCLostCallback
func ringbufCLostCallback(ctx unsafe.Pointer, cpu C.int, cnt C.ulonglong) {
	/* if verbose print? */
}
