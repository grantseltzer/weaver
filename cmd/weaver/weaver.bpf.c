#include "vmlinux.h"
#include <bpf/bpf_helpers.h>  
#include "weaver.h"     

char LICENSE[] SEC("license") = "GPL";

#define MAX_ENTRIES	7

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} output SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24); 
} dropped SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, MAX_ENTRIES);
	__type(key, u32);
	__type(value, struct parameter*);
} parameters SEC(".maps");

long ringbuffer_flags = 0;

SEC("uprobe/weaver")
int uprobe__weaver(struct pt_regs *ctx)
{
	u64 id = bpf_get_current_pid_tgid();
	u32 tgid = id >> 32;
	struct parameter *parameter;

    // Reserve space on the ringbuffer for the sample
	parameter = bpf_ringbuf_reserve(&output, sizeof(struct parameter), ringbuffer_flags);
	if (!parameter) {
		return 0;
    }

	void* stackAddr = (void*)ctx->sp;
	
	int idx;
	for (idx = 0; idx < 6; idx++) {
		struct parameter* param = (struct parameter*)bpf_map_lookup_elem(&parameters, (u32)idx);
		if (!parameter) {
			return 1;
		}

		u8 bytes[param->size];
		bpf_probe_read(bytes, param->size, param->start_offset);
		parameter->result = bytes;
		parameter->pid = tgid;
		bpf_ringbuf_submit(parameter, ringbuffer_flags);
	}
	
	return 0;
}
