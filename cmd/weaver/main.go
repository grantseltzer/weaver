package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/urfave/cli/v2"
)

var (
	globalDebug     bool
	globalDebugeBPF bool
	globalMode      modeOfOperation = PACKAGE_MODE

	globalOutput = os.Stdout
	globalError  = os.Stderr
)

func main() {

	app := &cli.App{
		Name:  "weaver",
		Usage: "Trace function executions within a specified Go binary file by any calling process",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "functions-file",
				Value:    "",
				Usage:    "specify a file which contains line sperated specifications of functions to trace. Each line is of the form: 'func-name(arg1_type, arg2_type, ...)'",
				Required: false,
				Aliases:  []string{"f"},
			},
			&cli.StringSliceFlag{
				Name:     "packages",
				Value:    cli.NewStringSlice("main"),
				Usage:    "specify a list of packages in the go binary to trace all of the functions in",
				Required: false,
				Aliases:  []string{"p"},
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
				Name:     "debug-ebpf",
				Value:    false,
				Usage:    "print eBPF program text before they're verified and loaded into the kernel",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "pid",
				Value:    -1,
				Usage:    "trace functions of already running Go binary using PID",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			return entry(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(globalError, err)
		fmt.Fprintln(globalError, "Try weaver --help")
		os.Exit(-1)
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

	// Turn on logging of eBPF programs
	if c.Bool("debug-ebpf") {
		globalDebugeBPF = true
	}

	// Just list acceptable golang types
	if c.Bool("types") {
		listAvailableTypes()
		return nil
	}

	if c.IsSet("functions-file") {
		globalMode = FUNC_FILE_MODE
	}

	var pid int
	var binaryArg string
	if c.IsSet("pid") {
		pid = c.Int("pid")
		binaryArg, err = getBinaryFromPID(pid)
		if err != nil {
			return err
		}
	} else {
		binaryArg = c.Args().Get(0)
		if binaryArg == "" {
			return errors.New("must specify a binary argument")
		}
	}

	binaryFullPath, err := filepath.Abs(binaryArg)
	if err != nil {
		return err
	}

	functionsFilePath := c.String("functions-file")

	// Read in functions file
	if globalMode == FUNC_FILE_MODE {
		contexts, err = readFunctionsFile(functionsFilePath)
		if err != nil {
			return err
		}
	} else {
		packagesToTrace := c.StringSlice("packages")
		contexts, err = read_symbols_from_binary(binaryFullPath, packagesToTrace)
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
