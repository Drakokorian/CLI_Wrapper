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
	return windows.SetPriorityClass(windows.Handle(cmd.Process.Pid), windows.BELOW_NORMAL_PRIORITY_CLASS)
}
