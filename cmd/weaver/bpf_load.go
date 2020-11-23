package main

/*
#cgo LDFLAGS: -lelf -lz

#include "weaver.skel.h"
#include "weaver.h"
#include "trace_helpers.h"

int libbpf_print_fn(enum libbpf_print_level level,
		    const char *format, va_list args)
{
	return vfprintf(stderr, format, args);
}

static int handle_event(void *ctx, void *data, size_t len)
{
	u64 *id = (u64*)data;
	printf("%d\n", *id);
	return 0;
}

int load_uprobe(char* binaryPath, size_t functionOffset) {

	int return_code = -1;
	struct weaver_bpf *obj;

	libbpf_set_print(libbpf_print_fn);

	obj = weaver_bpf__open_and_load();

	struct bpf_program *prog = bpf_object__find_program_by_name(obj->obj, "uprobe__weaver");
    if (!prog) {
		return_code = 1;
        goto cleanup;
	}

	struct bpf_link *link;
	link = bpf_program__atach_uprobe(prog, false, -1, binaryPath, functionOffset);
	if (!link) {
		return_code = 2;
		goto cleanup;
	}

	struct ring_buffer *ringbuffer;
	int ringbuffer_fd;
    ringbuffer_fd = bpf_map__fd(obj->maps.ringbuf);

	ringbuffer = ring_buffer__new(ringbuffer_fd, handle_event, NULL, NULL);
    if (!ringbuffer) {
		return_code =3;
		goto cleanup;
	}

	while (1) {
		// poll for new data with a timeout of -1 ms, waiting indefinitely
		ring_buffer__poll(ringbuffer, -1);
	}
cleanup:
	weaver_bpf__destroy(obj);
	return return_code;
}

*/
import "C"
import "fmt"

func load(target TraceTarget) error {
	ret := C.load_uprobe(C.CString(target.Name), C.size_t(target.Offset))
	return fmt.Errorf("uh oh %d", ret)
}
