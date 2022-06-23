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

SEC("uprobe/main.main")
int generic_function(struct pt_regs *ctx)
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

    bpf_printk("[] RAX: %d RBX: %d RCX: %d", out->registers.ax, out->registers.bx, out->registers.cx);
    bpf_printk("[] DI: %d SI: %d R8: %d", out->registers.di, out->registers.si, out->registers.r8);
    bpf_printk("[] R9: %d R10: %d R11: %d\n", out->registers.r9, out->registers.r10, out->registers.r11);

    bpf_ringbuf_submit(out, ringbuffer_flags);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
