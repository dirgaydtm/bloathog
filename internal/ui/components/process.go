package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/dirgaydtm/bloathog/internal/proc"
	"github.com/dirgaydtm/bloathog/internal/ui/theme"
)

// processItem wraps a process node and its visual prefix.
type processItem struct {
	node   proc.ProcessNode
	prefix string
}

func (p processItem) FilterValue() string { return p.node.Name }

// processDelegate renders process nodes in a list.
type processDelegate struct{}

// Pre-defined styles to save memory.
var (
	treeStyleName = lipgloss.NewStyle().Foreground(theme.ColorText)
	treeStyleRSS  = lipgloss.NewStyle().Foreground(theme.ColorAccent)
	treeStylePID  = lipgloss.NewStyle().Foreground(theme.ColorMuted)
)

func (d processDelegate) Height() int                             { return 2 }
func (d processDelegate) Spacing() int                            { return 0 }
func (d processDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d processDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	pi, ok := item.(processItem)
	if !ok {
		return
	}

	// Second line continues the tree's vertical lines under the first.
	indent2 := strings.NewReplacer("├─ ", "│  ", "└─ ", "   ").Replace(pi.prefix)
	pidStr := fmt.Sprintf("[%d]", pi.node.PID)

	nameStyle := treeStyleName
	rssStyle := treeStyleRSS
	if index == m.Index() {
		nameStyle = nameStyle.Background(theme.ColorAccent)
	} else if pi.node.Depth == 0 {
		nameStyle = nameStyle.Bold(true)
		rssStyle = rssStyle.Bold(true)
	}

	name := ansi.Truncate(pi.node.Name, m.Width()-lipgloss.Width(pi.prefix)-lipgloss.Width(pidStr)-2, "…")
	line1 := pi.prefix + nameStyle.Render(name) + " " + treeStylePID.Render(pidStr)
	line2 := fmt.Sprintf("%s%s  %s CPU", indent2,
		rssStyle.Render(fmt.Sprintf("%.1f MB", float64(pi.node.RSS)/(1024*1024))),
		treeStylePID.Render(fmt.Sprintf("%.1f%%", pi.node.CPU)))

	_, _ = fmt.Fprintf(w, "%s\n%s", line1, ansi.Truncate(line2, m.Width(), "…"))
}

// ProcessTreePanel displays the process tree UI.
type ProcessTreePanel struct {
	list    list.Model
	focused bool
}

// NewProcessTreePanel creates a new panel.
func NewProcessTreePanel(width, height int) ProcessTreePanel {
	l := list.New(nil, processDelegate{}, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Styles.NoItems = theme.StyleMuted
	return ProcessTreePanel{list: l}
}

// SetSize changes the panel dimensions.
func (p *ProcessTreePanel) SetSize(width, height int) { p.list.SetSize(width, height) }

// SetFocused toggles keyboard focus.
func (p *ProcessTreePanel) SetFocused(focused bool) { p.focused = focused }

// UpdateNodes refreshes the process list.
func (p *ProcessTreePanel) UpdateNodes(nodes []proc.ProcessNode) {
	items := make([]list.Item, len(nodes))
	for i, n := range nodes {
		prefix := ""
		if n.Depth > 0 && len(n.IsLast) >= n.Depth {
			for _, isLast := range n.IsLast[:n.Depth-1] {
				if isLast {
					prefix += "   "
				} else {
					prefix += "│  "
				}
			}
			if n.IsLast[n.Depth-1] {
				prefix += "└─ "
			} else {
				prefix += "├─ "
			}
		}
		items[i] = processItem{node: n, prefix: prefix}
	}
	p.list.SetItems(items)
}

// View draws the panel.
func (p *ProcessTreePanel) View() string {
	borderStyle := theme.StylePanel
	if p.focused {
		borderStyle = theme.StylePanelFocused
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		RenderTitle("Process Tree", p.list.Width()+4, p.focused),
		borderStyle.
			Width(p.list.Width()+2).
			Height(p.list.Height()).
			Render(p.list.View()),
	)
}

// ScrollUp moves selection up.
func (p *ProcessTreePanel) ScrollUp() { p.list.CursorUp() }

// ScrollDown moves selection down.
func (p *ProcessTreePanel) ScrollDown() { p.list.CursorDown() }
