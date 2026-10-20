package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dirgaydtm/bloathog/internal/monitor"
	"github.com/dirgaydtm/bloathog/internal/ring"
	"github.com/dirgaydtm/bloathog/internal/ui/components"
)

// quit triggers background shutdown without closing the UI.
func (m Model) quit() (tea.Model, tea.Cmd) {
	if !m.started {
		return m.finalizeQuit()
	}
	if m.killing {
		return m, nil
	}
	m.killing = true
	go monitor.KillProcess(m.cmd)
	if m.openStreams == 0 {
		return m, nil // waitCmd already running
	}
	// Don't wait for EOF: escaped descendants may hold the pipes open.
	return m, waitCmd(m.cmd)
}

// finalizeQuit runs after child process dies to safely close the UI.
func (m Model) finalizeQuit() (tea.Model, tea.Cmd) {
	m.quitting = true
	switch {
	case m.graph.Len() > 0:
		m.exitReport = components.RenderExitReport(m.stats, time.Since(m.startTime).Round(time.Second).String(), m.stats.PeakProcesses)
	case m.ExitCode == 0:
		m.exitReport = ring.Join(m.logs, "\n")
	}
	return m, tea.Sequence(m.header.Stop(), tea.Quit)
}
