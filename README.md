# Oster

<p align="center">
    <img src="DrManhattanGopher.png" alt="gopher" width="200"/>
</p>


Oster is a CLI tool that allows you to trace Go programs in order to inspect what values are passed to specified functions. It leverages eBPF attached to uprobes.


## Quick Start 

Take the following example program: 

```
package main

//go:noinline
func test_function(a int, d int32, e float64, r bool) {}

func main() {
	test_function(1, 2, 55555.111, true)
}
```

Let's say we want to know what values are passed to `test_function` whenever the program is run. Once the program is compiled we can call oster like this:

```
sudo oster --function-to-trace='main.test_function(int, int32, float64, bool)' ./test-prog
```

Notice that we have to specify the parameter data types. <i>(You can use `oster --types` to see what data types are supported.)</i>

Oster will sit idle without any output until `test-prog` is run and the `test_function` function is called. This will also work on a running Go Program.

```
[*] sudo ./bin/oster --function-to-trace='main.test_function(int, int32, float64, bool)' ./bin/test-prog&; sleep 1 && ./bin/test-prog
1
2                                                                                                                                       
55555.111000
true
```

## Note on supported types

Currently oster supports basic data types but getting support for user defined types is a high priority. Getting following types defined are also a work in progress:

- rune
- arbitrary pointers
- slices
- user/stdlib defined structs
- user/stdlib defined interfaces


## Dependencies

- [bcc](https://github.com/iovisor/bcc/blob/master/INSTALL.md)
- linux kernel version > 4.14 (please make bug reports if your kernel version doesn't work)

## Build

`make oster` will compile the oster binary to `bin/oster`

<i>Can't build? Please make an issue!</i>

## Roadmap

Check issues for tasks currently being tracked. Please open bug reports, i'm sure there are plenty :-)

Short term goals include:

- Testing
- Output options
- Tracing multiple functions at a time
- Deep pointer inspection
- Inspecting binaries for parameter data types instead of specifying them at the command line
- CI/CD infrastructre 

<i>image modified version of art by Ashley McNamara ([license](https://creativecommons.org/licenses/by-nc-sa/4.0/)) based on art by Renee French.</i>
