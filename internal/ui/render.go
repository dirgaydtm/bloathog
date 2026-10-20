package ui

import (
	"fmt"

	"github.com/dirgaydtm/bloathog/internal/ui/components"
	"github.com/dirgaydtm/bloathog/internal/ui/theme"
	"github.com/dirgaydtm/bloathog/internal/ui/types"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderLayout() string {
	if m.width == 0 {
		return ""
	}
	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal too small (min %dx%d)", minWidth, minHeight)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, theme.StyleWarning.Render(msg))
	}

	var activeKeys help.KeyMap
	switch m.focusTarget {
	case focusGraph:
		activeKeys = types.GraphKeyMap{KeyMap: m.keys}
	case focusLog:
		activeKeys = types.LogKeyMap{KeyMap: m.keys, InputMode: m.inputMode}
	default:
		activeKeys = types.ProcKeyMap{KeyMap: m.keys}
	}

	bottomSection := components.RenderHelpBar(m.helpModel, activeKeys, m.width)
	if m.inputMode {
		bottomSection = m.textInput.View() + "\n" + bottomSection
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.header.View(m.stats, m.inputMode, m.width),
		m.renderGraphAndProc(),
		m.logPanel.View(),
		bottomSection,
	)
}

func (m Model) renderGraphAndProc() string {
	graphW := m.width - m.procW
	isGraphFocused := m.focusTarget == focusGraph

	ramStyle, cpuStyle, graphData := theme.StyleAccent, theme.StyleMuted, m.graph
	if m.activeGraph == 1 {
		ramStyle, cpuStyle, graphData = cpuStyle, ramStyle, m.cpuGraph
	}
	graphTitle := ramStyle.Render("RAM Usage (MB)") + " │ " + cpuStyle.Render("CPU Usage (%)")

	panelStyle := theme.StylePanel
	if isGraphFocused {
		panelStyle = theme.StylePanelFocused
	}

	graphView := lipgloss.JoinVertical(lipgloss.Left,
		components.RenderTitle(graphTitle, graphW, isGraphFocused),
		panelStyle.
			Width(graphW-2). // Subtract 2 because lipgloss borders add 2 to total width
			Render(components.RenderGraph(graphData, graphW-4, m.graphHeight()-2, m.activeGraph)),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, m.procPanel.View(), graphView)
}
