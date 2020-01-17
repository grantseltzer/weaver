package main

//go:noinline
func test_function(a string, b byte, x int) {}

func main() {
	test_function("grant", byte(1), 7)
}
