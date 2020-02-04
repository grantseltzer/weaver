## Overview

Weaver is a useful debugging tool. The goal of weaver is to have [strace](https://linux.die.net/man/1/strace)-like functionality except for tracing functions in go programs.

Input:

- <b>A functions file</b> (using `-f` / `--functions-file`) -  A line delimited file containing specifications for function names and their argument types.

- <b>A compiled go program</b> - The actual binary which contains the functions specified in the functions file.

Behavior:

- Weaver reads the specified functions file. Each line is parsed into an internal data structure which represents the arguments.

- The stack offsets are then calculated based on [this logic](/docs/stack-offsets.md) and recorded in the same internal data structure.

-