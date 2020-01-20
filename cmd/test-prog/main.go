package main

//go:noinline
func test_function(int, string) {}

func main() {
	test_function(3, "hello")
}
