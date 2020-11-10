package bpfwrap

/*
#include "weaver.skel.h"

static struct env {
	pid_t pid;
	bool verbose;
} env = {};

void setPidFilter(int pid) {
	env.pid = (pid_t)pid;
}

void setVerboseFlag(bool verbose) {
	env.verbose = verbose;
}

int getRingBufFD() {
	return bpf_map__fd(obj->maps.ringbuf):
}

int libbpf_print_fn(enum libbpf_print_level level,
		    const char *format, va_list args)
{
	if (!env.verbose) {
		return 0;
    }
	return vfprintf(stderr, format, args);
}
*/
import "C"

import (
	"log"
)

func Init(targetPID int, verbose bool) {

	C.libbpf_set_print(C.libbpf_print_fn)

	obj := C.weaver_bpf__open()
	if obj == nil {
		log.Fatal("NIL OBJECT")
	}

	C.setPidFilter(targetPID)
	C.setVerboseFlag(verbose)

	cErr := C.weaver_bpf__load(obj)
	if cErr {
		log.Fatal("couldn't load obj")
	}

	cErr = C.weaver_bpf__attach(obj)
	if cErr {
		log.Fatal("couldn't attach obj")
	}

	ringbufFd := C.getRingBufFD()
	ringbuf = C.ring_buffer__new(ringbufFd, C.ringbufCallback, c.NULL, c.NULL)

	for {
		C.ring_buffer__poll(ringbuf, -1)
	}

	C.weaver_bpf__destroy(obj)
}
