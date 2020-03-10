package main

import (
	"fmt"
	"os"
)

func getBinaryFromPID(pid int) (string, error) {
	// Weaver should have the rights to read /proc
	exe := fmt.Sprintf("/proc/%d/exe", pid)
	_, err := os.Stat(exe)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%s: binary does not exist: %s", err.Error(), exe)
		}

		return "", err
	}

	return exe, nil
}
