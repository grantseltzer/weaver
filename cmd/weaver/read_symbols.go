package main

import (
	"debug/elf"
	"strings"
)

func read_symbols_from_binary(binaryPath string, packagesToTrace []string) ([]functionTraceContext, error) {
	elfFile, err := elf.Open(binaryPath)
	if err != nil {
		return nil, err
	}

	syms, err := elfFile.Symbols()
	if err != nil {
		return nil, err
	}

	contexts := []functionTraceContext{}

	for _, sym := range syms {

		ok := validateSymbol(sym)
		if !ok {
			continue
		}

		for _, pkg := range packagesToTrace {
			if strings.HasPrefix(sym.Name, pkg+".") {
				x := functionTraceContext{
					Filters:      filters{},
					binaryName:   binaryPath,
					HasArguments: false,
					FunctionName: sym.Name,
				}
				contexts = append(contexts, x)
			}
		}
	}

	return contexts, nil
}

func validateSymbol(sym elf.Symbol) bool {

	return strings.Count(sym.Name, ".") == 1
}
