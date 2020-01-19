package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"

	bpf "github.com/iovisor/gobpf/bcc"
)

const textTemplate = `
	#include <uapi/linux/ptrace.h>

	BPF_PERF_OUTPUT(events);

	inline int print_stack(struct pt_regs *ctx) {
		void* stackAddr = (void*)ctx->sp;

		int i;
		char y;
		for (i = 0; i < 50; i++) {
			char *x = (char *)stackAddr+i;
			bpf_probe_read(&y, sizeof(y), x);
			events.perf_submit(ctx, &y, sizeof(y));
		}

		return 0;
	}
`

func main() {
	bpfModule := bpf.NewModule(textTemplate, []string{})
	defer bpfModule.Close()

	uprobeFd, err := bpfModule.LoadUprobe("print_stack")
	if err != nil {
		log.Fatal(err)
	}

	err = bpfModule.AttachUprobe(os.Args[1], os.Args[2], uprobeFd, -1)
	if err != nil {
		log.Fatalf("could not attach uprobe to symbol: %s: %s", "test_function", err.Error())
	}

	// Set up bpf perf map to use for output
	table := bpf.NewTable(bpfModule.TableId("events"), bpfModule)
	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		i := 0
		for {
			value := <-channel
			x1 := binary.LittleEndian.Uint32(value)
			fmt.Printf("SP+%d:\t%d (0x%x)\n", i, x1, x1)
			i++
		}
	}()

	perfMap.Start()
	<-sig
	perfMap.Stop()
}
