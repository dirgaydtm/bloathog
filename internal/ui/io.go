package ui

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaydtm/bloathog/internal/monitor"
)

// logBatchMsg carries a batch of logs.
type logBatchMsg struct {
	msgs []monitor.LogMsg
	ch   <-chan monitor.LogMsg
}

type streamDoneMsg struct{}

// readLinesCmd streams lines from a reader into a channel.
func readLinesCmd(r io.Reader, isStderr bool) tea.Cmd {
	ch := make(chan monitor.LogMsg, 256)
	go func() {
		defer close(ch)
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadString('\n')
			if line != "" {
				line = strings.TrimRight(line, "\r\n")
				// Keep only what a terminal would show after the last \r.
				ch <- monitor.LogMsg{Line: line[strings.LastIndexByte(line, '\r')+1:], IsStderr: isStderr}
			}
			if err != nil {
				return
			}
		}
	}()
	return drainCmd(ch)
}

// drainCmd reads up to 100 log lines at once for performance.
func drainCmd(ch <-chan monitor.LogMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return streamDoneMsg{}
		}

		batch := []monitor.LogMsg{msg}
		for len(batch) < 100 {
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
// Run after both streams hit EOF: Wait closes the pipes.
func waitCmd(cmd *exec.Cmd) tea.Cmd {
	return func() tea.Msg {
		var exitErr *exec.ExitError
		if errors.As(cmd.Wait(), &exitErr) {
			return monitor.ChildExitMsg{ExitCode: exitErr.ExitCode()}
		}
		return monitor.ChildExitMsg{}
	}
}
