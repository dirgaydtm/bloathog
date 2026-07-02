package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
const (
	ColorAccent  = lipgloss.Color("#0A84FF") // System Blue
	ColorSuccess = lipgloss.Color("#32D74B") // System Green
	ColorWarning = lipgloss.Color("#FFD60A") // System Yellow
	ColorDanger  = lipgloss.Color("#FF453A") // System Red
	ColorMuted   = lipgloss.Color("#8E8E93") // System Gray
	ColorText    = lipgloss.Color("#F2F2F7") // System Gray 6 (White)
	ColorSubtext = lipgloss.Color("#98989D") // Secondary Text
	ColorBorder  = lipgloss.Color("#38383A") // Very subtle border
	ColorBg      = lipgloss.Color("#1C1C1E") // macOS Dark Bg
	ColorBgPanel = lipgloss.Color("#2C2C2E") // macOS Elevated Bg
)

// Text styles
var (
	StyleBold    = lipgloss.NewStyle().Bold(true)
	StyleMuted   = lipgloss.NewStyle().Foreground(ColorMuted)
	StyleSubtext = lipgloss.NewStyle().Foreground(ColorSubtext)

	StyleAccent  = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	StyleSuccess = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	StyleWarning = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	StyleDanger  = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)

	StyleHeader = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true).
			Padding(0, 1)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Bold(false)

	StyleValue = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)
)

// Panel border style
var (
	BorderStyle = lipgloss.RoundedBorder()

	TabTitleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)

	StylePanel = lipgloss.NewStyle().
			Border(BorderStyle).
			BorderTop(false).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	StylePanelFocused = lipgloss.NewStyle().
				Border(BorderStyle).
				BorderTop(false).
				BorderForeground(ColorAccent).
				Padding(0, 1)

	StylePanelTitle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1)
)