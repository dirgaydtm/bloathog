package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/dirgaa/bloathog/internal/ui/theme"
	"github.com/dirgaa/bloathog/internal/ui/types"
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
		sp:      sp,
		sw:      stopwatch.NewWithInterval(time.Second),
		CmdStr:  cmdStr,
	}
}

func (m HeaderModel) Init() tea.Cmd {
	return m.sp.Tick
}

func (m *HeaderModel) SetStarted(started bool) tea.Cmd {
	m.started = started
	if started {
		return m.sw.Start()
	}
	return nil
}

func (m *HeaderModel) Stop() tea.Cmd {
	return m.sw.Stop()
}

func (m *HeaderModel) Update(msg tea.Msg) (HeaderModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if !m.started {
			var cmd tea.Cmd
			m.sp, cmd = m.sp.Update(msg)
			cmds = append(cmds, cmd)
		}
	case stopwatch.TickMsg, stopwatch.StartStopMsg:
		var cmd tea.Cmd
		m.sw, cmd = m.sw.Update(msg)
		cmds = append(cmds, cmd)
	}

	return *m, tea.Batch(cmds...)
}

func (m HeaderModel) View(stats types.MonitorState, width int) string {
	var status string
	if !m.started {
		status = m.sp.View() + " starting…"
	} else {
		status = theme.StyleSuccess.Render("●") + " running  " +
			theme.StyleMuted.Render("⏱ "+m.sw.View())
	}

	statsStr := fmt.Sprintf("RAM: %.1fMB (Peak: %.1fMB) │ CPU: %.1f%% (Peak: %.1f%%)",
		stats.CurrentRSSMB(), stats.PeakRSSMB(),
		stats.CurrentCPU, stats.PeakCPU)

	left := theme.StyleAccent.Render("bloathog") + "  " + theme.StyleSuccess.Render(m.CmdStr)
	
	leftW := lipgloss.Width(left)
	statusW := lipgloss.Width(status)
	
	contentW := width - 2 
	availableW := contentW - leftW - statusW - 4 
	if availableW < 0 {
		availableW = 0
	}

	centerTxt := statsStr
	if lipgloss.Width(centerTxt) > availableW {
		if availableW > 3 {
			runes := []rune(statsStr)
			if len(runes) > availableW {
				centerTxt = string(runes[:availableW-3]) + "..."
			}
		} else {
			centerTxt = ""
		}
	}
	
	centerStyled := theme.StyleLabel.Render(centerTxt)

	spacerL := (contentW - leftW - statusW - lipgloss.Width(centerStyled)) / 2
	spacerR := contentW - leftW - statusW - lipgloss.Width(centerStyled) - spacerL
	if spacerL < 0 { spacerL = 0 }
	if spacerR < 0 { spacerR = 0 }

	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		Render(lipgloss.JoinHorizontal(lipgloss.Top,
			left,
			lipgloss.NewStyle().Width(spacerL).Render(""),
			centerStyled,
			lipgloss.NewStyle().Width(spacerR).Render(""),
			status,
		))
}
