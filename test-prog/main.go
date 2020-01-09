package main

//go:noinline
func test_function(a int, d int32, e float64, r bool) {}

func main() {
	test_function(1, 2, 55555.111, true)
}
