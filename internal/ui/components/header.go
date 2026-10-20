package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/dirgaydtm/bloathog/internal/ui/theme"
	"github.com/dirgaydtm/bloathog/internal/ui/types"
)

type HeaderModel struct {
	sp      spinner.Model
	sw      stopwatch.Model
	started bool
	CmdStr  string
}

func NewHeaderModel(cmdStr string) HeaderModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(theme.ColorAccent)

	return HeaderModel{
		sp:     sp,
		sw:     stopwatch.NewWithInterval(time.Second),
		CmdStr: cmdStr,
	}
}

func (m HeaderModel) Init() tea.Cmd {
	return m.sp.Tick
}

func (m *HeaderModel) Start() tea.Cmd {
	m.started = true
	return m.sw.Start()
}

func (m *HeaderModel) Stop() tea.Cmd {
	return m.sw.Stop()
}

func (m *HeaderModel) Update(msg tea.Msg) (HeaderModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case spinner.TickMsg:
		if !m.started {
			m.sp, cmd = m.sp.Update(msg)
		}
	case stopwatch.TickMsg, stopwatch.StartStopMsg:
		m.sw, cmd = m.sw.Update(msg)
	}
	return *m, cmd
}

func (m HeaderModel) View(stats types.MonitorState, inputMode bool, width int) string {
	var status string
	timer := theme.StyleMuted.Render("⏱ " + m.sw.View())
	if inputMode {
		status = theme.StyleWarning.Render("●") + " input    " + timer
	} else if !m.started {
		status = m.sp.View() + " starting…"
	} else {
		status = theme.StyleSuccess.Render("●") + " running  " + timer
	}

	statsStr := fmt.Sprintf("RAM: %.1fMB (Peak: %.1fMB) │ CPU: %.1f%% (Peak: %.1f%%)",
		stats.CurrentRSSMB(), stats.PeakRSSMB(),
		stats.CurrentCPU, stats.PeakCPU)

	statusW := lipgloss.Width(status)
	contentW := width - 2

	cmdStr := ansi.Truncate(m.CmdStr, contentW-statusW-11, "…") // 11 = "bloathog  " + 1 gap
	left := theme.StyleAccent.Render("bloathog") + "  " + theme.StyleSuccess.Render(cmdStr)
	leftW := lipgloss.Width(left)

	center := theme.StyleLabel.Render(ansi.Truncate(statsStr, contentW-leftW-statusW-4, "..."))

	gap := max(0, contentW-leftW-statusW-lipgloss.Width(center))
	row := left + strings.Repeat(" ", gap/2) + center + strings.Repeat(" ", gap-gap/2) + status
	return lipgloss.NewStyle().Width(width).Padding(0, 1).Render(row)
}
