package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	globalDebug bool
)

func main() {

	app := &cli.App{
		Name:  "weaver",
		Usage: "Trace specific executions in Go programs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "functions-file",
				Value:    "",
				Usage:    "specify a file which contains line sperated specifications of functions to trace. Each line is of the form: 'func-name(arg1_type, arg2_type, ...)'. Use `weaver --types` for a list of accepted types.",
				Required: false,
				Aliases:  []string{"f"},
			},
			&cli.BoolFlag{
				Name:     "types",
				Value:    false,
				Usage:    "list accepted types for function parameters",
				Required: false,
				Aliases:  []string{"t"},
			},
			&cli.BoolFlag{
				Name:     "debug",
				Value:    false,
				Usage:    "turn on debug logging to stderr",
				Required: false,
				Aliases:  []string{"d"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.NumFlags() == 0 {
				fmt.Println("Use weaver --help")
				os.Exit(0)
			}
			return entry(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func entry(c *cli.Context) error {

	var (
		contexts []functionTraceContext
		err      error
	)

	// Turn on debug logigng
	if c.Bool("debug") {
		globalDebug = true
	}

	// Just list acceptable golang types
	if c.Bool("types") {
		listAvailableTypes()
		return nil
	}

	binaryFullPath, err := filepath.Abs(c.Args().Get(0))
	if err != nil {
		return err
	}

	functionsFilePath := c.String("functions-file")

	// Read in functions file
	if functionsFilePath != "" {
		contexts, err = read_functions_file(functionsFilePath)
		if err != nil {
			return err
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	runtimeContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Install eBPF program for each function to trace
	for i := range contexts {

		contexts[i].binaryName = binaryFullPath

		// Load uprobe and BPF code. This will block until Ctrl-C or an error occurs.
		go loadUprobeAndBPFModule(&contexts[i], runtimeContext)
	}

	<-sig

	return nil
}
