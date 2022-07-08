package main

//go:noinline
func test_single_int(x int) {}

func main() {
	test_single_int(7)
}
