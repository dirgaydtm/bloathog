package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dirgaa/bloathog/internal/monitor"
	"github.com/dirgaa/bloathog/internal/ring"
	"github.com/dirgaa/bloathog/internal/ui/components"
)

// quit triggers background shutdown without closing the UI.
func (m Model) quit() (tea.Model, tea.Cmd) {
	if !m.started {
		return m.finalizeQuit()
	}
	if !m.quitting {
		m.quitting = true
		go monitor.KillProcess(m.cmd)
	}
	return m, nil
}

// finalizeQuit runs after child process dies to safely close the UI.
func (m Model) finalizeQuit() (tea.Model, tea.Cmd) {
	m.quitting = true
	if m.graph.Len() == 0 {
		if m.ExitCode == 0 {
			m.exitReport = ring.Join(m.logs, "\n")
		} else {
			m.exitReport = ""
		}
	} else {
		m.exitReport = components.RenderExitReport(m.stats, time.Since(m.startTime).Round(time.Second).String(), m.peakPrcs)
	}
	return m, tea.Sequence(m.header.Stop(), tea.Quit)
}
