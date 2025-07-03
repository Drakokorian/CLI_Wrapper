//go:build windows

package app

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os/exec"
)

func throttleProcess(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	handle, err := windows.OpenProcess(windows.PROCESS_SET_INFORMATION, false, uint32(cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("open process: %w", err)
	}
	defer windows.CloseHandle(handle)
	return windows.SetPriorityClass(handle, windows.BELOW_NORMAL_PRIORITY_CLASS)
}
