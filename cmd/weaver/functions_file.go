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
	contexts := make([]functionTraceContext, len(functionStringsToTrace))

	for i, funcString := range functionStringsToTrace {

		if funcString == "" {
			continue
		}

		debugLog("parsing: %s\n", funcString)

		err := parseFunctionAndArgumentTypes(&contexts[i], funcString)
		if err != nil {
			return nil, fmt.Errorf("could not parse function string '%s': %s", funcString, err.Error())
		}

		err = determineStackOffsets(&contexts[i])
		if err != nil {
			return nil, fmt.Errorf("could not determine stack offsets of arguments: %s", err.Error())
		}

	}

	return contexts, nil
}
