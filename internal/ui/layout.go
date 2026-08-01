package ui

func (m *Model) relayout() {
	m.procW = max(25, min(m.width/2, int(float64(m.width)*0.20)))
	m.procPanel.SetSize(m.procW-4, m.graphHeight()-2)
	m.logPanel.SetSize(m.width-4, m.logPanelHeight()-2)
}

func (m Model) graphHeight() int {
	return max(7, min(19, int(float64(m.height-2)*0.6)))
}

func (m Model) logPanelHeight() int {
	reserved := 2 + m.graphHeight()
	if m.inputMode {
		reserved++
	}
	return max(3, m.height-reserved)
}
