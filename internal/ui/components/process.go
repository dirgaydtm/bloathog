package components

import (
	"github.com/dirgaa/bloathog/internal/ui/theme"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/dirgaa/bloathog/internal/proc"
)

// processItem wraps proc.ProcessNode and its pre-calculated lipgloss tree prefix.
type processItem struct {
	node   proc.ProcessNode
	prefix string
}

func (p processItem) FilterValue() string { return p.node.Name }

// processDelegate is a custom list.ItemDelegate for process nodes.
type processDelegate struct{}

// Package-level styles to avoid allocating new lipgloss.Style objects on every render call.
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

	indent := pi.prefix
	
	// Create second line indent: replace branching chars with vertical or space chars
	indent2 := strings.ReplaceAll(indent, "├─ ", "│  ")
	indent2 = strings.ReplaceAll(indent2, "└─ ", "   ")

	rssMB := float64(pi.node.RSS) / (1024 * 1024)
	rssStr := fmt.Sprintf("%.1f MB", rssMB)
	pidStr := fmt.Sprintf("[%d]", pi.node.PID)
	cpuStr := fmt.Sprintf("%.1f%%", pi.node.CPU)

	nameStyle := treeStyleName
	rssStyle := treeStyleRSS

	if index == m.Index() {
		nameStyle = nameStyle.Foreground(theme.ColorText).Background(theme.ColorAccent)
	} else if pi.node.Depth == 0 {
		nameStyle = nameStyle.Bold(true)
		rssStyle = rssStyle.Bold(true)
	}

	name := pi.node.Name
	// Optional: we don't need intense truncation logic here since it's on its own line,
	// but we'll apply it just in case name is ridiculously long.
	maxNameW := m.Width() - lipgloss.Width(indent) - lipgloss.Width(pidStr) - 2
	if lipgloss.Width(name) > maxNameW && maxNameW > 3 {
		runes := []rune(name)
		if len(runes) > maxNameW {
			name = string(runes[:maxNameW-1]) + "…"
		}
	}

	line1 := fmt.Sprintf("%s%s %s", indent, nameStyle.Render(name), treeStylePID.Render(pidStr))
	line2 := fmt.Sprintf("%s%s  %s CPU", indent2, rssStyle.Render(rssStr), treeStylePID.Render(cpuStr))
	
	fmt.Fprintf(w, "%s\n%s", line1, line2)
}

// ProcessTreePanel is a bubbles/list panel showing the live process tree.
type ProcessTreePanel struct {
	list    list.Model
	focused bool
}

// NewProcessTreePanel creates a process tree panel with given dimensions.
func NewProcessTreePanel(width, height int) ProcessTreePanel {
	l := list.New(nil, processDelegate{}, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.NoItems = theme.StyleMuted
	return ProcessTreePanel{list: l}
}

// SetSize resizes the panel.
func (p *ProcessTreePanel) SetSize(width, height int) { p.list.SetSize(width, height) }

// SetFocused sets whether this panel has keyboard focus.
func (p *ProcessTreePanel) SetFocused(focused bool) { p.focused = focused }

// UpdateNodes replaces the displayed process nodes.
func (p *ProcessTreePanel) UpdateNodes(nodes []proc.ProcessNode) {
	if len(nodes) == 0 {
		p.list.SetItems(nil)
		return
	}

	// Use lipgloss/tree to pre-calculate the perfect indentation strings for each node.
	// We use "." as a placeholder for the actual content.
	root := tree.Root(".")
	stack := []*tree.Tree{root}

	for i := 1; i < len(nodes); i++ {
		n := nodes[i]
		if n.Depth <= len(stack) && n.Depth > 0 {
			stack = stack[:n.Depth]
		}
		parent := stack[len(stack)-1]
		child := tree.Root(".")
		parent.Child(child)
		stack = append(stack, child)
	}

	lines := strings.Split(root.String(), "\n")

	items := make([]list.Item, len(nodes))
	for i, n := range nodes {
		prefix := ""
		if i < len(lines) {
			prefix = strings.TrimSuffix(lines[i], ".")
			// Compress the tree indentation to save horizontal space
			prefix = strings.ReplaceAll(prefix, "├── ", "├─ ")
			prefix = strings.ReplaceAll(prefix, "└── ", "└─ ")
			prefix = strings.ReplaceAll(prefix, "│   ", "│  ")
			prefix = strings.ReplaceAll(prefix, "    ", "   ")
		}
		items[i] = processItem{node: n, prefix: prefix}
	}
	p.list.SetItems(items)
}

// View renders the panel with title and border.
func (p *ProcessTreePanel) View() string {
	header := RenderTitle("Process Tree", p.list.Width()+4, p.focused)

	borderStyle := theme.StylePanel
	if p.focused {
		borderStyle = theme.StylePanelFocused
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		borderStyle.
			Width(p.list.Width()+2).
			Height(p.list.Height()).
			Render(p.list.View()),
	)
}

// ScrollUp moves the cursor up.
func (p *ProcessTreePanel) ScrollUp(n int) {
	for i := 0; i < n; i++ {
		p.list.CursorUp()
	}
}

// ScrollDown moves the cursor down.
func (p *ProcessTreePanel) ScrollDown(n int) {
	for i := 0; i < n; i++ {
		p.list.CursorDown()
	}
}

// Update passes a tea.Msg to the underlying list (useful for mouse clicks).
func (p *ProcessTreePanel) Update(msg tea.Msg) { p.list, _ = p.list.Update(msg) }
