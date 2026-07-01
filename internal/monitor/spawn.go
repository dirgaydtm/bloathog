package monitor

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

// SpawnCmd returns a tea.Cmd that spawns the child process and returns an
// InternalStartMsg so the UI model can wire goroutines on receipt.
func SpawnCmd(command string, args []string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(command, args...)

		// Unix: run in its own process group to kill the whole tree later
		if runtime.GOOS != "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		if err := cmd.Start(); err != nil {
			return ErrorMsg{Err: err}
		}

		return InternalStartMsg{
			RootPID: int32(cmd.Process.Pid),
			Cmd:     cmd,
			Stdout:  stdout,
			Stderr:  stderr,
		}
	}
}

// KillProcess terminates the child process and its entire tree.
func KillProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	if runtime.GOOS == "windows" {
		killer := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", cmd.Process.Pid))
		_ = killer.Run()
		return
	}

	// Unix: kill the entire process group (negative PID = PGID)
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
		time.Sleep(200 * time.Millisecond)
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	} else {
		_ = cmd.Process.Signal(os.Interrupt)
	}
}
