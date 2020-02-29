package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/urfave/cli/v2"
)

var (
	globalDebug bool
	globalJSON  bool
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
			&cli.BoolFlag{
				Name:     "json",
				Value:    false,
				Usage:    "toggle for output to be in JSON format",
				Required: false,
				Aliases:  []string{"JSON", "j"},
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

	// Turn on debug logigng
	if c.Bool("json") {
		globalJSON = true
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

	// wg is used to communicate back with the main thread when
	// all uprobe/eBPFs are installed
	var wg sync.WaitGroup

	// Install eBPF program for each function to trace
	for i := range contexts {
		wg.Add(1)
		contexts[i].binaryName = binaryFullPath
		go loadUprobeAndBPFModule(&contexts[i], runtimeContext, &wg)
	}

	go func() {
		wg.Wait()
		debugLog("All probes installed\n")
	}()

	<-sig

	return nil
}
