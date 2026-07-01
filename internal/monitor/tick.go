package monitor

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dirgaa/bloathog/internal/proc"
)

const tickInterval = 500 * time.Millisecond

// TickCmd returns a tea.Cmd that collects process tree stats periodically.
func TickCmd(rootPID int32) tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		stats := proc.CollectTree(rootPID)
		return TickMsg{Stats: stats}
	})
}
