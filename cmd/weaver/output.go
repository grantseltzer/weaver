package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/olekukonko/tablewriter"
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

var mutex = &sync.Mutex{}

func printOutput(o output) error {

	if globalJSON {
		b, err := json.Marshal(o)
		if err != nil {
			return fmt.Errorf("could not marshal output to JSON: %s", err.Error())
		}
		fmt.Println(string(b))
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function Name", "Arg Position", "Type", "Value", "Proc Name", "PID", "PPID"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for i, arg := range o.Args {
		line := []string{o.FunctionName, fmt.Sprintf("%d", i), arg.Type, arg.Value, o.ProcInfo.Comm,
			fmt.Sprintf("%d", o.ProcInfo.Pid), fmt.Sprintf("%d", o.ProcInfo.Ppid)}
		table.Append(line)
	}

	// tablewriter doesn't support asynchronous renders
	mutex.Lock()
	table.Render()
	mutex.Unlock()

	return nil
}

func debugLog(format string, a ...interface{}) {
	if globalDebug {
		fmt.Fprintf(os.Stderr, "\x1b[96m"+format, a...)
	}
}
