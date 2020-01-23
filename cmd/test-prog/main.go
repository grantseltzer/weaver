package main

//go:noinline
func test_function(int, [2]int) {}

func main() {
	test_function(3, [2]int{1, 2})
}
