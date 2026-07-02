package ui

import (
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/spinner"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/monitor"
	"github.com/dirgaa/bloathog/internal/ui/components"
	"github.com/dirgaa/bloathog/internal/ui/types"
)

const (
	maxGraphSamples  = 120
	maxLogLines      = 5000
	focusLog         = 0
	focusProc        = 1
	focusGraph       = 2
	focusCount       = 3
)

// Model is the root Bubble Tea model for bloathog.
type Model struct {
	projectInfo types.ProjectInfo

	// Process state
	rootPID  int32
	cmd      *exec.Cmd
	started  bool
	quitting bool
	exitCode int

	// Monitor state (streaming stats)
	stats    types.MonitorState
	peakPrcs int

	// Dynamic buffers
	graph       []float64
	cpuGraph    []float64
	activeGraph int // 0 = RAM, 1 = CPU
	logs        []string
	logDirty bool

	// UI components
	logPanel    components.LogPanel
	procPanel   components.ProcessTreePanel
	header      components.HeaderModel
	helpModel   help.Model
	keys        types.KeyMap
	focusTarget int

	// Layout
	width    int
	height   int
	procW    int // cached panel width, set in relayout()
	showHelp bool

	// Exit state
	exitReport string
	startTime  time.Time
	fatalErr   error
}

func NewModel(info types.ProjectInfo) tea.Model {
	h := help.New()
	h.ShowAll = false

	cmdStr := info.Command
	for _, a := range info.Args {
		cmdStr += " " + a
	}

	m := Model{
		projectInfo: info,
		header:      components.NewHeaderModel(cmdStr),
		helpModel:   h,
		keys:        types.DefaultKeyMap(),
		focusTarget: focusLog,
		logPanel:    components.NewLogPanel(80, 10),
		procPanel:   components.NewProcessTreePanel(40, 10),
	}

	m.logPanel.SetFocused(true)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.header.Init(), monitor.SpawnCmd(m.projectInfo.Command, m.projectInfo.Args))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.relayout()

	case tea.KeyMsg:
		// Context-aware graph switching
		if m.focusTarget == focusGraph && key.Matches(msg, m.keys.SwitchGraph) {
			m.activeGraph = 1 - m.activeGraph
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m.quit()
		case key.Matches(msg, m.keys.Tab):
			m.focusTarget = (m.focusTarget + 1) % focusCount
			m.logPanel.SetFocused(m.focusTarget == focusLog)
			m.procPanel.SetFocused(m.focusTarget == focusProc)
		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.helpModel.ShowAll = m.showHelp
		case key.Matches(msg, m.keys.Up):
			switch m.focusTarget {
			case focusLog:
				m.logPanel.ScrollUp(1)
			case focusProc:
				m.procPanel.ScrollUp(1)
			}
		case key.Matches(msg, m.keys.Down):
			switch m.focusTarget {
			case focusLog:
				m.logPanel.ScrollDown(1)
			case focusProc:
				m.procPanel.ScrollDown(1)
			}
		}

	case monitor.InternalStartMsg:
		m.cmd = msg.Cmd
		m.rootPID = msg.RootPID
		m.started = true
		m.startTime = time.Now()
		cmds = append(cmds,
			m.header.SetStarted(true),
			monitor.TickCmd(msg.RootPID),
			readLinesCmd(msg.Stdout, false),
			readLinesCmd(msg.Stderr, true),
			waitCmd(msg.Cmd),
		)

	case monitor.TickMsg:
		if !m.quitting {
			rss := msg.Stats.TotalRSS
			m.stats.CurrentRSS = rss
			if rss > m.stats.PeakRSS {
				m.stats.PeakRSS = rss
			}
			m.stats.RunningSumRSS += float64(rss)
			m.stats.SampleCount++
			m.stats.ActiveProcesses = msg.Stats.ProcessCount
			if msg.Stats.ProcessCount > m.peakPrcs {
				m.peakPrcs = msg.Stats.ProcessCount
			}
			m.graph = append(m.graph, float64(rss)/(1024*1024))
			if len(m.graph) > maxGraphSamples {
				m.graph = m.graph[1:]
			}
			
			cpu := msg.Stats.TotalCPU
			m.stats.CurrentCPU = cpu
			if cpu > m.stats.PeakCPU {
				m.stats.PeakCPU = cpu
			}
			m.stats.RunningSumCPU += cpu
			m.cpuGraph = append(m.cpuGraph, cpu)
			if len(m.cpuGraph) > maxGraphSamples {
				m.cpuGraph = m.cpuGraph[1:]
			}

			m.procPanel.UpdateNodes(msg.Stats.Nodes)
			cmds = append(cmds, monitor.TickCmd(m.rootPID))
		}

	case nextLineWithChanMsg:
		m.logs = append(m.logs, formatLogLine(msg.msg))
		if len(m.logs) > maxLogLines {
			m.logs = m.logs[1:]
		}
		m.logDirty = true
		cmds = append(cmds, drainCmd(msg.ch))

	case monitor.ChildExitMsg:
		m.exitCode = msg.ExitCode
		return m.quit()

	case monitor.ErrorMsg:
		m.fatalErr = msg.Err
		return m.quit()

	case spinner.TickMsg, stopwatch.TickMsg, stopwatch.StartStopMsg:
		var cmd tea.Cmd
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.logDirty {
		m.logPanel.UpdateContent(m.logs)
		m.logDirty = false
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.exitReport != "" {
		return m.exitReport
	}
	return m.renderLayout()
}
