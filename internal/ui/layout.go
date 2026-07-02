package ui

func (m *Model) relayout() {
	procW := int(float64(m.width) * 0.20)
	if procW < 25 {
		procW = 25
	}	
	if procW > m.width/2 {
		procW = m.width / 2
	}
	
	m.procW = procW
	
	gh := m.graphHeight()
	m.procPanel.SetSize(procW-4, gh-2)

	// components.LogPanel header=1, border=1. Total=2.
	m.logPanel.SetSize(m.width-4, m.logPanelHeight()-2)
}

func (m Model) graphHeight() int {
	available := m.height - 2 // reserved for header (1) + help (1)
	
	target := 19
	maxAllowed := int(float64(available) * 0.6) // max 60% of available height
	if target > maxAllowed {
		target = maxAllowed
	}
	
	if target < 7 { // absolute minimum height for graph visibility
		target = 7
	}
	
	return target
}

func (m Model) logPanelHeight() int {
	reserved := 2                // header (1 line) + help (1 line)
	reserved += m.graphHeight()  // graph panel with border
	h := m.height - reserved
	if h < 3 {
		h = 3
	}
	return h
}
