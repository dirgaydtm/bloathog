package ui

import (
	"github.com/dirgaa/bloathog/internal/ui/components"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/monitor"
)

func (m *Model) quit() (tea.Model, tea.Cmd) {
	m.quitting = true
	monitor.KillProcess(m.cmd)
	m.exitReport = components.RenderExitReport(m.stats, formatDuration(time.Since(m.startTime)), m.peakPrcs)
	return m, tea.Sequence(m.header.Stop(), tea.Quit)
}

// formatDuration converts a duration to a human-readable string like "2h 03m 15s".
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
