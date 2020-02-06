## Overview

Weaver is a useful debugging tool. The goal of weaver is to have [strace](https://linux.die.net/man/1/strace)-like functionality except for tracing functions in go programs.

Input:

- <b>A functions file</b> (using `-f` / `--functions-file`) -  A line delimited file containing specifications for function names and their argument types.

- <b>A compiled go program</b> - The actual binary which contains the functions specified in the functions file.

Behavior:

- Weaver reads the specified functions file. Each line is parsed into an internal data structure (trace context) which represents the arguments.

- The stack offsets are then calculated based on [this logic](/docs/stack-offsets.md) and recorded in the same trace context structure.

- The trace context for each function is compiled into eBPF programs using a text template.

- Each eBPF program is loaded into the kernel attached to their corresponding uprobes.

- Weaver listens on a single perf buffer for argument values from the eBPF programs once their triggered by running the corresponding program.

- Arguments are outputted in the configured manner.