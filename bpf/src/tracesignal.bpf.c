#include "vmlinux.h"
#include <bpf/bpf_helpers.h>  
#include "tracesignal.h"     

char LICENSE[] SEC("license") = "GPL";

const volatile pid_t target_tgid = 0;

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} ringbuf SEC(".maps");

long ringbuffer_flags = 0;

static __always_inline
bool trace_allowed(u32 tgid)
{
	if (target_tgid && target_tgid != tgid) {
		return false;
    }
	return true;
}

SEC("tracepoint/signal/signal_generate")
int tracepoint__signal_signal_generate(struct trace_event_raw_signal_generate* ctx)
{
    u64 id = bpf_get_current_pid_tgid();
	u32 tgid = id >> 32;
	struct process_info *process;

    if (!trace_allowed(tgid)) {
        return 1;
    }

    // Reserve space on the ringbuffer for the sample
	process = bpf_ringbuf_reserve(&ringbuf, sizeof(*process), ringbuffer_flags);
	if (!process) {
		return 0;
    }

    process->pid = tgid;
    process->signal = ctx->sig;
    bpf_probe_read_kernel(process->comm, sizeof(process->comm), ctx->comm);
  
    bpf_ringbuf_submit(process, ringbuffer_flags);

    return 0;
}

// SEC("tracepoint/signal/signal_deliver")
// int tracepoint__signal_signal_deliver(struct trace_event_raw_signal_deliver* ctx)
// {

// }
