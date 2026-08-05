//go:build windows

package monitor

import (
	"fmt"
	"os/exec"
)

func setProcessGroup(cmd *exec.Cmd) {
}

func killProcessGroup(cmd *exec.Cmd) {
	killer := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", cmd.Process.Pid))
	_ = killer.Run()
}
