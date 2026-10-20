package components

import (
	"github.com/dirgaydtm/bloathog/internal/ring"
	"github.com/dirgaydtm/bloathog/internal/ui/theme"

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
	atBottom := lp.vp.AtBottom()
	lp.vp.Width = width
	lp.vp.Height = height
	if atBottom {
		lp.vp.GotoBottom()
	}
}

// SetFocused sets whether this panel has keyboard focus.
func (lp *LogPanel) SetFocused(focused bool) { lp.focused = focused }

func (lp *LogPanel) UpdateContent(buffer *ring.Buffer[string]) {
	atBottom := lp.vp.AtBottom()
	lp.vp.SetContent(ring.Join(buffer, "\n"))
	if atBottom {
		lp.vp.GotoBottom()
	}
}

// ScrollUp scrolls the viewport up by one line.
func (lp *LogPanel) ScrollUp() { lp.vp.ScrollUp(1) }

// ScrollDown scrolls the viewport down by one line.
func (lp *LogPanel) ScrollDown() { lp.vp.ScrollDown(1) }

// View renders the log panel with a title and border.
func (lp *LogPanel) View() string {
	borderStyle := theme.StylePanel
	if lp.focused {
		borderStyle = theme.StylePanelFocused
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		RenderTitle("System Logs", lp.vp.Width+4, lp.focused),
		borderStyle.
			Width(lp.vp.Width+2).
			Height(lp.vp.Height).
			Render(lp.vp.View()),
	)
}
