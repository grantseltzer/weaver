//+build ignore
#include "vmlinux.h"
#include <string.h>
#include <bpf/bpf_helpers.h>  

#ifdef asm_inline
#undef asm_inline
#define asm_inline asm
#endif

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
} events SEC(".maps");

long ringbuffer_flags = 0;

struct out_context {
    char stack[50];
    struct pt_regs registers;
};

SEC("uprobe/dump")
int dump(struct pt_regs *ctx)
{
    struct out_context *out;
    out = bpf_ringbuf_reserve(&events, sizeof(struct out_context), ringbuffer_flags);
    if (!out) {
        return 1;
    }

    struct out_context tmp;

    void* stackAddr = (void*)ctx->sp;
    int i;
    char y;

    for (i = 0; i < 50; i++) {
        bpf_probe_read(tmp.stack+i, sizeof(char), stackAddr+i);
    }

    tmp.registers = *ctx;
    bpf_probe_read(out, sizeof(tmp), &tmp);
    bpf_ringbuf_submit(out, ringbuffer_flags);
    return 0;
}

SEC("uprobe/repeat")
int repeat(struct pt_regs *ctx)
{
    // Read the address of the calling routine (return address) from the top of the stack
    void* stackAddr = (void*)ctx->sp;
    unsigned long returnAddress;
    bpf_probe_read((void*)&returnAddress, sizeof(returnAddress), stackAddr);
   
    // Overwrite the return address with the top of the current routine
    returnAddress = ctx->ip;
    bpf_probe_write_user(stackAddr, (void*)&returnAddress, sizeof(returnAddress));
    
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
