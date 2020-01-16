package main

import (
	"fmt"
	"strings"
)

//go:noinline
func test_function(a string, x int) {
	fmt.Println(&a, a)
}

func main() {
	variable := strings.Split("AAAAAAAAAAAA", ",")
	d := variable[0]
	fmt.Println(&d)
	test_function(d, 7)
}
