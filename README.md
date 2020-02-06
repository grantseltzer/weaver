# Weaver

<p align="center">
    <img src="DrManhattanGopher.png" alt="gopher" width="200"/>
</p>


Weaver is a CLI tool that allows you to trace Go programs in order to inspect what values are passed to specified functions. It leverages eBPF attached to uprobes.

[![Go Report Card](https://goreportcard.com/badge/github.com/grantseltzer/weaver)](https://goreportcard.com/report/github.com/grantseltzer/weaver)


## Quick Start 

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
sudo ./bin/weaver -f ./cmd/test-prog/test_functions_file.txt ./bin/test-prog | jq
```


Weaver will then sit idle without any output until `test-prog` is run and the `test_function` function is called. This will also work on a running Go Program.

```json
{
  "FunctionName": "main.other_test_function",
  "Args": [
    {
      "Type": "RUNE",
      "Value": "a"
    },
    {
      "Type": "INT64",
      "Value": "33"
    }
  ]
}
{
  "FunctionName": "main.test_function",
  "Args": [
    {
      "Type": "INT",
      "Value": "3"
    },
    {
      "Type": "INT_ARRAY",
      "Value": "1, 2"
    }
  ]
}

```

## Note on supported types

Currently weaver supports basic data types but getting support for user defined types is a high priority. Getting following types defined are also a work in progress:

- arbitrary pointers
- slices
- user/stdlib defined structs
- user/stdlib defined interfaces


## Dependencies

- [bcc](https://github.com/iovisor/bcc/blob/master/INSTALL.md)
- linux kernel version > 4.14 (please make bug reports if your kernel version doesn't work)

## Build

`make weaver` will compile the weaver binary to `bin/weaver`

<i>Can't build? Please make an issue!</i>

## Roadmap

Check issues for tasks currently being tracked. Please open bug reports, i'm sure there are plenty :-)

Short term goals include:

- Testing
- Output options
- Deep pointer inspection
- Inspecting binaries for parameter data types instead of specifying them at the command line
- CI/CD infrastructre 

<i>image modified version of art by Ashley McNamara ([license](https://creativecommons.org/licenses/by-nc-sa/4.0/)) based on art by Renee French.</i>
