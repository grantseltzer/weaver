package main

//go:noinline
func test_function(int, [2]int) {}

//go:noinline
func other_test_function(rune, int64) {}

func main() {
	test_function(3, [2]int{1, 2})
	other_test_function('a', 33)
}
