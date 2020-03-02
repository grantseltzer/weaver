package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type output struct {
	FunctionName string      `json:"functionName"`
	Args         []outputArg `json:"args"`
	ProcInfo     procInfo    `json:"procInfo"`
}

type outputArg struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func printOutput(o output) error {

	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("could not marshal output to JSON: %s", err.Error())
	}
	fmt.Println(string(b))

	return nil
}

func debugLog(format string, a ...interface{}) {
	if globalDebug {
		fmt.Fprintf(os.Stderr, "\x1b[96m"+format+"\x1b[0m", a...)
	}
}

func debugeBPFLog(format string, a ...interface{}) {
	if globalDebugeBPF {
		fmt.Fprintf(os.Stderr, "\x1b[96m"+format+"\x1b[0m", a...)
	}
}
