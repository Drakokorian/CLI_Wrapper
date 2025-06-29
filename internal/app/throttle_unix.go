//go:build !windows

package app

import (
	"fmt"
	"os/exec"
	"syscall"
)

func throttleProcess(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	return syscall.Setpriority(syscall.PRIO_PROCESS, cmd.Process.Pid, 10)
}
