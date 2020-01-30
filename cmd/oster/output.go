package main

import (
	"encoding/json"
	"fmt"
)

type output struct {
	FunctionName string
	Args         []outputArg
}

type outputArg struct {
	Type  string
	Value string
}

func printOutput(o output) error {
	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("could not marshal output to JSON: %s", err.Error())
	}
	fmt.Println(string(b))
	return nil
}
