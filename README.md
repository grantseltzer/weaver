# Weaver

<b>PLEASE READ!</b> - This project is inactive. I'm working on this exactly functionality for Datadog as part of the [Dynamic Instrumentation](https://docs.datadoghq.com/dynamic_instrumentation/) product. Check out the datadog agent for code there!

<p align="center">
    <img src="DrManhattanGopher.png" alt="gopher" width="200"/>
</p>


Weaver is a CLI tool that allows you to trace Go programs in order to inspect what values are passed to specified functions. It leverages eBPF attached to uprobes.

[![Go Report Card](https://goreportcard.com/badge/github.com/grantseltzer/weaver)](https://goreportcard.com/report/github.com/grantseltzer/weaver)


## Quick Start 

There are two modes of operation, one that uses a 'functions file', and one that extracts a symbol table from a passed binary and filters by Go packages. More information on functionality in [docs](/docs).

### Functions file

Take the following example program: 

<i>test_prog.go</i>
```go
package main

//go:noinline
func test_function(int, [2]int) {}

//go:noinline
func other_test_function(rune, int64) {}

func main() {
	test_function(3, [2]int{1, 2})
	other_test_function('a', 33)
}
```

Let's say we want to know what values are passed to `test_function` and `other_test_function` whenever the program is run. Once the program is compiled (`make`) we just have to create a file which specifies each function to trace:

<i>functions_to_trace.txt</i>
```
main.test_function(int, [2]int)
main.other_test_function(rune, int64)
```

Notice that we have to specify the parameter data types. <i>(You can use `weaver --types` to see what data types are supported.)</i>

Now we can call `weaver` like so:

```
sudo weaver -f /path/to/functions_to_trace.txt /path/to/test-prog-binary
```

Weaver will then sit idle without any output until `test-prog` is run and the `test_function` and `other_test_function` functions are called. This will also work on an already running Go Program.

```
{"functionName":"main.other_test_function","args":[{"type":"RUNE","value":"a"},{"type":"INT64","value":"33"}],"procInfo":{"pid":43300,"ppid":42754,"comm":"test-prog-binar"}}
{"functionName":"main.test_function","args":[{"type":"INT","value":"3"},{"type":"INT_ARRAY","value":"1, 2"}],"procInfo":{"pid":43300,"ppid":42754,"comm":"test-prog-binar"}}
```

### Package mode

For the same example Go program as above, you can choose to not specify a functions file. The command would like like this:

```
sudo weaver /path/to/test-prog-binary
```

This will default to only tracing functions in the `main` package, however you can use the `--packages` flag to specify a comma seperated list of packages (typially of the form `github.com/x/y`)

Output does include argument vlaues in this mode.

```
{"functionName":"main.main","procInfo":{"pid":44411,"ppid":42754,"comm":"test-prog-binar"}}
{"functionName":"main.test_function","procInfo":{"pid":44411,"ppid":42754,"comm":"test-prog-binar"}}
```

## Note on supported types

Currently weaver supports basic data types but getting support for user defined types is a high priority. Getting following types defined are a work in progress:

- user/stdlib defined structs
- user/stdlib defined interfaces


## System Dependencies

- [bcc](https://github.com/iovisor/bcc/blob/master/INSTALL.md) / bcc-devel
- linux kernel version > 4.14 (please make bug reports if your kernel version doesn't work)

## Build

`make` will compile the weaver binary to `bin/weaver` (It also creates the smoke test binary and print-stack utility)

<i>Can't build? Please make an issue!</i>

## Roadmap

Check issues for tasks currently being tracked. Please open bug reports, i'm sure there are plenty :-)

Short term goals include:

- Testing
- Output options
- Inspecting binaries for parameter data types instead of specifying them with a functions file
- CI/CD infrastructre 

<i>image modified version of art by Ashley McNamara ([license](https://creativecommons.org/licenses/by-nc-sa/4.0/)) based on art by Renee French.</i>
