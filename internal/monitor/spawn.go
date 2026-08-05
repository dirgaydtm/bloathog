package monitor

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// SpawnCmd returns a tea.Cmd that spawns the child process and returns an
// InternalStartMsg so the UI model can wire goroutines on receipt.
func SpawnCmd(command string, args []string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(command, args...)

		setProcessGroup(cmd)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		// Inject FORCE_COLOR=1 to trick dev servers into outputting colors
		cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "CLICOLOR_FORCE=1")

		if err := cmd.Start(); err != nil {
			return ErrorMsg{Err: err}
		}

		return InternalStartMsg{
			RootPID: int32(cmd.Process.Pid),
			Cmd:     cmd,
			Stdin:   stdin,
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

	killProcessGroup(cmd)
}
