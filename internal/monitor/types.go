package monitor

import (
	"io"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/proc"
)

// InternalStartMsg is sent when the child process has been successfully spawned.
// It carries the running exec.Cmd so the model can wire goroutines and kill on quit.
type InternalStartMsg struct {
	RootPID int32
	Cmd     *exec.Cmd
	Stdout  io.Reader
	Stderr  io.Reader
}

// TickMsg is sent every 500ms with updated process tree stats.
type TickMsg struct {
	Stats proc.TreeStats
}

// LogMsg is sent when a new line is captured from child stdout/stderr.
type LogMsg struct {
	Line     string
	IsStderr bool
}

// ChildExitMsg is sent when the child process terminates on its own.
type ChildExitMsg struct {
	ExitCode int
}

// ErrorMsg carries a fatal error to display before quitting.
type ErrorMsg struct {
	Err error
}

// Cmd is a helper that wraps a tea.Cmd result — unused externally,
// kept for internal convenience.
type Cmd = tea.Cmd
