package main

/*
#include "weaver.skel.h"

int load_uprobe() {

	struct weaver_bpf *obj;

	obj = weaver_bpf__open();
	if (!obj) {
		fprintf(stderr, "failed to open and/or load BPF object\n");
		return -1; // FIXME: how does returns work? have constants for display from Go?
	}

	err = uprobetest_bpf__load(obj);
    if (err) {
		fprintf(stderr, "failed to load BPF object: %d\n", err);
		goto cleanup;
    }


}

int libbpf_print_fn(enum libbpf_print_level level,
		    const char *format, va_list args)
{
	// TODO: send to error ringbuf
	return vfprintf(stderr, format, args);
}

static int handle_event(void *ctx, void *data, size_t len)
{
	// TODO: send to output ringbuf
    struct process_info *s = (struct process_info*)data;
	printf("%d >%d<\n", s->pid, s->arg);
	return 0;
}

void handle_lost_events(void *ctx, int cpu, __u64 lost_cnt)
{
	// TODO: send to lost events ringbuf
	fprintf(stderr, "Lost %llu events on CPU #%d!\n", lost_cnt, cpu);
}

*/
import "C"

func load(target TraceTarget) error {
	C.load_uprobe()
}
