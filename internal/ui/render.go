package ui

import (
	"github.com/dirgaa/bloathog/internal/ui/components"
	"github.com/dirgaa/bloathog/internal/ui/theme"
	"github.com/dirgaa/bloathog/internal/ui/types"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderLayout() string {
	var activeKeys help.KeyMap
	if m.focusTarget == focusGraph {
		activeKeys = types.GraphKeyMap{KeyMap: m.keys}
	} else {
		activeKeys = types.ScrollKeyMap{KeyMap: m.keys}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.header.View(m.stats, m.width),
		m.renderGraphAndProc(),
		m.logPanel.View(),
		components.RenderHelpBar(m.helpModel, activeKeys, m.width),
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
	var graphData []float64
	
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
