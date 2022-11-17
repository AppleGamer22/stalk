// https://gist.github.com/hallazzang/76f3970bfc949831808bbebc8ca15209
package main

import (
	"os"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

type windowsProcess struct {
	Pid    int
	Handle uintptr
}

type ProcessExitGroup windows.Handle

func NewProcessExitGroup() (ProcessExitGroup, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return ProcessExitGroup(handle), nil
}

func (g ProcessExitGroup) Dispose() error {
	return windows.CloseHandle(windows.Handle(g))
}

func (g ProcessExitGroup) AddProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(
		windows.Handle(g),
		windows.Handle((*windowsProcess)(unsafe.Pointer(p)).Handle))
}

func start(command string, arguments ...string) error {
	processMutex.Lock()
	defer processMutex.Unlock()
	if process != nil {
		if err := kill(false); err != nil {
			return err
		}
	}

	group, err := NewProcessExitGroup()
	if err != nil {
		return err
	}
	defer group.Dispose()

	process = exec.Command(command, arguments...)
	process.Stdout = os.Stdout
	process.Stdin = os.Stdin
	process.Stderr = os.Stderr

	if err := group.AddProcess(process.Process); err != nil {
		return err
	}

	return process.Start()
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
