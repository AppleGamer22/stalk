package main

import (
	"os"
	"os/exec"
	"sync"
)

var (
	process      *exec.Cmd
	processMutex sync.Mutex
)

func start(command string, arguments ...string) error {
	processMutex.Lock()
	defer processMutex.Unlock()
	if process != nil {
		if err := kill(false); err != nil {
			return err
		}
	}

	process = exec.Command(command, arguments...)
	process.Stdout = os.Stdout
	process.Stdin = os.Stdin
	process.Stderr = os.Stderr

	if err := process.Start(); err != nil {
		return err
	}

	return nil
}

func kill(lock bool) error {
	if lock {
		processMutex.Lock()
		defer processMutex.Unlock()
	}
	if process == nil {
		return nil
	}
	if err := process.Process.Kill(); err != nil {
		return err
	}
	_, err := process.Process.Wait()
	return err
}
