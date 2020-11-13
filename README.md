# Weaver REFACTOR!

Welcome to the <b>refactor</b> branch!

I'm currently working through refactoring weaver to use libbpf. There are very important motivations for this and there are some significant challenges to this. I try to highlight them below:

## Motivations

- BCC, and especially the Go bcc bindings, are poorly maintained. As a result of this we don't get important features like ringbuffers (ordering of events out of kernel) and all future features that we'll want to use.
- The Go bcc library also has significant bugs in it that aren't getting fixed such as broken tail call support.
- libbpf is maintained within the kernel which will always be a well maintained project.
- libbpf will always have the latest and greatest features of bpf.
- One of the great things about using libbpf is the ability to do CO-RE which enables easy portability.

## Challenges

- No major linux distributions are shipping with CO-RE support (Ubuntu 21.04 will be first I believe)
- Lack of general documentation around libbpf for uprobes (This is the first project i've seen doing this at all to be honest)
- Requirement of using CGO to attach the C user space bpf code to the bpf code itself

# Other refactoring 

- Use ringbuf instead of a perf buffer for ordered events
- Remove the need for a functions file by properly parsing dwarfinfo (see cmd/weaver/parse_dwarf.go)
- More!