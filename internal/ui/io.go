package ui

import (
	"bufio"
	"io"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/monitor"
	"github.com/dirgaa/bloathog/internal/ui/theme"
)

// logBatchMsg carries a batch of logs.
type logBatchMsg struct {
	msgs []monitor.LogMsg
	ch   <-chan monitor.LogMsg
}

// formatLogLine adds a warning mark to stderr.
func formatLogLine(msg monitor.LogMsg) string {
	if msg.IsStderr {
		return theme.StyleWarning.Render("!") + " " + msg.Line
	}
	return "  " + msg.Line
}

// readLinesCmd streams lines from a reader into a channel.
func readLinesCmd(r io.Reader, isStderr bool) tea.Cmd {
	ch := make(chan monitor.LogMsg, 256)
	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 64*1024), 64*1024)
		for scanner.Scan() {
			ch <- monitor.LogMsg{Line: scanner.Text(), IsStderr: isStderr}
		}
		if err := scanner.Err(); err != nil {
			ch <- monitor.LogMsg{Line: "[bloathog warning] error reading stream: " + err.Error(), IsStderr: true}
		}
		close(ch)
	}()
	return drainCmd(ch)
}

// drainCmd reads up to 100 log lines at once for performance.
func drainCmd(ch <-chan monitor.LogMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}

		batch := []monitor.LogMsg{msg}
		for i := 0; i < 99; i++ {
			select {
			case m, ok := <-ch:
				if !ok {
					return logBatchMsg{msgs: batch, ch: ch}
				}
				batch = append(batch, m)
			default:
				return logBatchMsg{msgs: batch, ch: ch}
			}
		}
		return logBatchMsg{msgs: batch, ch: ch}
	}
}

// waitCmd waits for child process exit.
func waitCmd(cmd *exec.Cmd) tea.Cmd {
	return func() tea.Msg {
		err := cmd.Wait()
		code := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				code = exitErr.ExitCode()
			}
		}
		return monitor.ChildExitMsg{ExitCode: code}
	}
}
