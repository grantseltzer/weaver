package main

//go:noinline
func test_function(x rune) {}

func main() {
	test_function('a')
}
