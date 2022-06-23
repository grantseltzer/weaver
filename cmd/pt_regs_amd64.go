package main

type trace_context struct {
	MemoryStack Mem_stack
	Registers   Pt_regs
}

type Mem_stack [50]byte

type Pt_regs struct {
	R15     int64
	R14     int64
	R13     int64
	R12     int64
	Bp      int64
	Bx      int64
	R11     int64
	R10     int64
	R9      int64
	R8      int64
	Ax      int64
	Cx      int64
	Dx      int64
	Si      int64
	Di      int64
	Orig_ax int64
	Ip      int64
	Cs      int64
	Flags   int64
	Sp      int64
	Ss      int64
}
