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
	u64 *idptr;
    // Reserve space on the ringbuffer for the sample
	idptr = bpf_ringbuf_reserve(&output, sizeof(id, ringbuffer_flags), ringbuffer_flags);
	if (!idptr) {
		return 0;
    }
	
	bpf_ringbuf_submit(idptr, ringbuffer_flags);

	return 0;
}
