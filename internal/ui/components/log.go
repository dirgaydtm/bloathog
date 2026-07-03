package components

import (
	"github.com/dirgaa/bloathog/internal/ui/theme"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// LogPanel wraps bubbles/viewport with follow-mode and ring-buffer rendering.
type LogPanel struct {
	vp      viewport.Model
	focused bool
}

// NewLogPanel creates a new LogPanel with the given dimensions.
func NewLogPanel(width, height int) LogPanel {
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle()
	return LogPanel{vp: vp}
}

// SetSize resizes the viewport.
func (lp *LogPanel) SetSize(width, height int) {
	lp.vp.Width = width
	lp.vp.Height = height
}

// SetFocused sets whether this panel has keyboard focus.
func (lp *LogPanel) SetFocused(focused bool) { lp.focused = focused }

func (lp *LogPanel) UpdateContent(lines []string) {
	isAtBottom := lp.vp.AtBottom()
	content := strings.Join(lines, "\n")
	lp.vp.SetContent(content)
	if isAtBottom || len(lines) == 1 {
		lp.vp.GotoBottom()
	}
}

// ScrollUp scrolls the viewport up by one line.
func (lp *LogPanel) ScrollUp(n int) { lp.vp.LineUp(n) }

// ScrollDown scrolls the viewport down by one line.
func (lp *LogPanel) ScrollDown(n int) { lp.vp.LineDown(n) }

// View renders the log panel with a title and border.
func (lp *LogPanel) View() string {
	header := RenderTitle("System Logs", lp.vp.Width+4, lp.focused)

	borderStyle := theme.StylePanel
	if lp.focused {
		borderStyle = theme.StylePanelFocused
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		borderStyle.
			Width(lp.vp.Width+2).
			Height(lp.vp.Height).
			Render(lp.vp.View()),
	)
}
