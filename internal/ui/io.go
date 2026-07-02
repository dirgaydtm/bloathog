package ui

import (
	"bufio"
	"io"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/monitor"
	"github.com/dirgaa/bloathog/internal/ui/theme"
)

// formatLogLine prefixes stderr lines with a warning indicator.
func formatLogLine(msg monitor.LogMsg) string {
	if msg.IsStderr {
		return theme.StyleWarning.Render("!") + " " + msg.Line
	}
	return "  " + msg.Line
}

// readLinesCmd starts a goroutine that reads all lines from r into a channel,
// then returns a drainCmd to pull them into the tea event loop one by one.
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

// drainCmd reads one LogMsg from the channel and schedules itself to continue.
func drainCmd(ch <-chan monitor.LogMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil // channel closed, stop draining
		}
		// Return the log message; the Update handler will call drainCmd again
		// via nextLineMsg
		return nextLineWithChanMsg{msg: msg, ch: ch}
	}
}

// waitCmd waits for the child process to exit and sends ChildExitMsg.
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

// nextLineWithChanMsg carries a log line AND the channel to continue draining.
type nextLineWithChanMsg struct {
	msg monitor.LogMsg
	ch  <-chan monitor.LogMsg
}
