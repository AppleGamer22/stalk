//go:build unix

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func start(command string, arguments ...string) (*exec.Cmd, error) {
	process := exec.Command(command, arguments...)
	process.Stdout = os.Stdout
	// process.Stdin = os.Stdin
	process.Stderr = os.Stderr
	process.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return process, process.Start()
}

func kill(process *exec.Cmd) error {
	if process == nil {
		return nil
	}

	if err := syscall.Kill(-process.Process.Pid, syscall.SIGKILL); err != nil {
		return err
	}

	_, err := process.Process.Wait()
	return err
}
