package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dirgaa/bloathog/internal/ui/theme"
)

// RenderTitle creates a single-line embedded header that connects to a BorderTop(false) panel.
func RenderTitle(title string, width int, isFocused bool) string {
	borderColor := theme.ColorBorder
	titleColor := theme.ColorSubtext
	if isFocused {
		borderColor = theme.ColorAccent
		titleColor = theme.ColorAccent
	}

	titleText := theme.TabTitleStyle.
		Foreground(titleColor).
		Render(title)

	// Width of "╭─" is 2.
	lineLen := width - lipgloss.Width(titleText) - 2
	if lineLen < 1 {
		lineLen = 1
	}

	leftCorner := lipgloss.NewStyle().Foreground(borderColor).Render("╭─")
	line := lipgloss.NewStyle().Foreground(borderColor).Render(strings.Repeat("─", lineLen-1) + "╮")

	return leftCorner + titleText + line
}
