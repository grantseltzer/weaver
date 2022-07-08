package main

import (
	"fmt"
	"time"
)

//go:noinline
func test_print_something() {
	time.Sleep(time.Second)
	fmt.Println("It's February 2nd!")
}

func main() {
	test_print_something()
}
