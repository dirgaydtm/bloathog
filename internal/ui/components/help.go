package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/dirgaydtm/bloathog/internal/ui/theme"
)

// renderHelpBar renders the help footer using bubbles/help.
func RenderHelpBar(h help.Model, keys help.KeyMap, width int) string {
	h.Styles.ShortKey = lipgloss.NewStyle().
		Background(theme.ColorAccent).
		Foreground(theme.ColorText).
		Padding(0, 1)

	h.ShortSeparator = "   "

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(theme.ColorText).
		Padding(0, 1).
		Render(ansi.Truncate(h.View(keys), width-2, "…"))
}
