package main

import (
	"errors"
	"os"
	"syscall"
)

func getProcessRunningStatus(pid int) (*os.Process, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	//double check if process is running and alive
	//by sending a signal 0
	//NOTE : syscall.Signal is not available in Windows

	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return proc, nil
	}

	if err == syscall.ESRCH {
		return nil, errors.New("process not running")
	}

	// default
	return nil, errors.New("process running but query operation not permitted")
}

func getBinaryFromPID(pid int) (string, error) {
	// Weaver should have the rights to read /proc

	// Check if process is running
	_, err := getProcessRunningStatus(pid)
	if err != nil {
		return "", err
	}

	// Get executable name

	return "", nil
}
