package ui

import (
	"github.com/dirgaydtm/bloathog/internal/ring"
	"github.com/dirgaydtm/bloathog/internal/ui/components"
	"github.com/dirgaydtm/bloathog/internal/ui/theme"
	"github.com/dirgaydtm/bloathog/internal/ui/types"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderLayout() string {
	var activeKeys help.KeyMap
	if m.focusTarget == focusGraph {
		activeKeys = types.GraphKeyMap{KeyMap: m.keys}
	} else if m.focusTarget == focusLog {
		activeKeys = types.LogKeyMap{KeyMap: m.keys, InputMode: m.inputMode}
	} else {
		activeKeys = types.ProcKeyMap{KeyMap: m.keys}
	}

	var bottomSection string
	if m.inputMode {
		bottomSection = lipgloss.JoinVertical(lipgloss.Left,
			m.textInput.View(),
			components.RenderHelpBar(m.helpModel, activeKeys, m.width),
		)
	} else {
		bottomSection = components.RenderHelpBar(m.helpModel, activeKeys, m.width)
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
	if graphW < 12 {
		graphW = 12
	}

	leftCol := m.procPanel.View()

	isGraphFocused := m.focusTarget == focusGraph

	var graphHeaderTitle string
	var graphData *ring.Buffer[float64]

	ramTitle := "RAM Usage (MB)"
	cpuTitle := "CPU Usage (%)"

	if m.activeGraph == 0 {
		graphHeaderTitle = theme.StyleAccent.Render(ramTitle) + " │ " + theme.StyleMuted.Render(cpuTitle)
		graphData = m.graph
	} else {
		graphHeaderTitle = theme.StyleMuted.Render(ramTitle) + " │ " + theme.StyleAccent.Render(cpuTitle)
		graphData = m.cpuGraph
	}

	graphHeader := components.RenderTitle(graphHeaderTitle, graphW, isGraphFocused)

	panelStyle := theme.StylePanel
	if isGraphFocused {
		panelStyle = theme.StylePanelFocused
	}

	graphView := lipgloss.JoinVertical(lipgloss.Left,
		graphHeader,
		panelStyle.
			Width(graphW-2). // Subtract 2 because lipgloss borders add 2 to total width
			Render(components.RenderGraph(graphData, graphW-4, m.graphHeight()-2, m.activeGraph)),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, graphView)
}
