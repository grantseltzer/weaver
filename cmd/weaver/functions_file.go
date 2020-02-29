package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func read_functions_file(path string) ([]functionTraceContext, error) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read functions file")
	}

	functionStringsToTrace := strings.Split(string(content), "\n")

	contexts := []functionTraceContext{}

	for _, funcString := range functionStringsToTrace {

		if funcString == "" || funcString == "\n" {
			continue
		}

		debugLog("parsing: %s\n", funcString)

		newContext := functionTraceContext{}

		err := parseFunctionAndArgumentTypes(&newContext, funcString)
		if err != nil {
			return nil, fmt.Errorf("could not parse function string '%s': %s", funcString, err.Error())
		}

		err = determineStackOffsets(&newContext)
		if err != nil {
			return nil, fmt.Errorf("could not determine stack offsets of arguments: %s", err.Error())
		}

		contexts = append(contexts, newContext)

	}

	return contexts, nil
}
