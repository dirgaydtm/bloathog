package components

import (
	"github.com/dirgaa/bloathog/internal/ui/theme"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// renderHelpBar renders the help footer using bubbles/help.
func RenderHelpBar(h help.Model, keys help.KeyMap, width int) string {
	h.Styles.ShortKey = lipgloss.NewStyle().
		Background(theme.ColorAccent).
		Foreground(theme.ColorText).
		Padding(0, 1)

	h.ShortSeparator = "   "

	bar := h.View(keys)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(theme.ColorText).
		Padding(0, 1).
		Render(bar)
}
