package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	globalDebug         bool
	globalDebugeBPF     bool
	globalPackageFilter bool
	globalOutput        = os.Stdout
	globalError         = os.Stderr
)

func main() {

	app := &cli.App{
		Name:  "weaver",
		Usage: "Trace function executions within a specified Go binary file by any calling process",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "packages",
				Value:    cli.NewStringSlice("main"),
				Usage:    "specify a list of packages in the go binary to trace all of the functions in",
				Required: false,
				Aliases:  []string{"p"},
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

	// Turn on debug logigng
	if c.Bool("debug") {
		globalDebug = true
	}

	// Turn on logging of eBPF programs
	if c.Bool("debug-ebpf") {
		globalDebugeBPF = true
	}

	path := c.Args().Get(0)
	if path == "" {
		return errors.New("must specify a binary argument")
	}

	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	binaryFullPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	packagesToTrace := c.StringSlice("packages")

	filters := TraceFilter{
		packages: packagesToTrace,
	}

	traceTargets, err := GetTargets(binaryFullPath, filters)
	if err != nil {
		return err
	}

	s, _ := json.Marshal(traceTargets)
	fmt.Printf("%s", s)

	return nil
}
