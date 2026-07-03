package components

import (
	"fmt"
	"strings"

	"github.com/dirgaa/bloathog/internal/ui/theme"
	"github.com/dirgaa/bloathog/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	reportCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorMuted).
		Padding(0, 1).
		Margin(1, 1)

	reportTitleStyle = lipgloss.NewStyle().
		Foreground(theme.ColorText).
		Bold(true).
		Align(lipgloss.Center).
		Width(32)

	reportLabelStyle = lipgloss.NewStyle().
		Foreground(theme.ColorSubtext).
		Width(16).
		Align(lipgloss.Left)

	reportValueStyle = lipgloss.NewStyle().
		Foreground(theme.ColorText).
		Bold(true).
		Width(16).
		Align(lipgloss.Right)
)

// RenderExitReport renders the final summary box shown when bloathog quits.
func RenderExitReport(state types.MonitorState, durationStr string, peakProcesses int) string {
	rows := []struct{ label, value string }{
		{"Duration", durationStr},
		{"Peak RAM", fmt.Sprintf("%.2f MB", state.PeakRSSMB())},
		{"Average RAM", fmt.Sprintf("%.2f MB", state.AverageRSSMB())},
		{"Peak CPU", fmt.Sprintf("%.2f %%", state.PeakCPU)},
		{"Average CPU", fmt.Sprintf("%.2f %%", state.AverageCPU())},
		{"Peak Processes", fmt.Sprintf("%d", peakProcesses)},
		{"Total Samples", fmt.Sprintf("%d", state.SampleCount)},
	}

	var sb strings.Builder
	
	sb.WriteString(reportTitleStyle.Render("Session Summary"))
	sb.WriteString("\n\n")
	
	for i, row := range rows {
		label := reportLabelStyle.Render(row.label)
		value := reportValueStyle.Render(row.value)
		
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, label, value))
		if i < len(rows)-1 {
			sb.WriteString("\n")
		}
	}

	return reportCardStyle.Render(sb.String())
}