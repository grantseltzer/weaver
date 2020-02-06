package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
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

	if globalJSON {
		b, err := json.Marshal(o)
		if err != nil {
			return fmt.Errorf("could not marshal output to JSON: %s", err.Error())
		}
		fmt.Println(string(b))
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function Name", "Arg Position", "Type", "Value"})

	for i, arg := range o.Args {
		line := []string{o.FunctionName, fmt.Sprintf("%d", i), arg.Type, arg.Value}
		table.Append(line)
	}
	table.Render()

	return nil
}
