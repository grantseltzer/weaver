package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var globalDebug bool

func main() {

	app := &cli.App{
		Name:  "oster",
		Usage: "Trace specific executions in Go programs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "function-to-trace",
				Value:    "",
				Usage:    "specify function name and arguments of the form: 'func-name(arg1_type, arg2_type, ...)'. Use oster --types for list of accepted types.",
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
				fmt.Println("Use oster --help")
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

	if c.Bool("debug") {
		globalDebug = true
	}

	// Just list acceptable golang types
	if c.Bool("types") {
		listAvailableTypes()
		return nil
	}

	// Initialize tracing info
	var context traceContext

	fullPath, err := filepath.Abs(c.Args().Get(0))
	if err != nil {
		log.Fatal(err)
	}

	context.binaryName = fullPath

	err = parseFunctionAndArgumentTypes(&context, c.String("function-to-trace"))
	if err != nil {
		log.Fatal(err)
	}

	err = determineStackOffsets(&context)
	if err != nil {
		log.Fatalf("could not determine stack offsets of arguments: %s", err.Error())
	}

	// Load uprobe
	err = createBPFModule(&context)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
