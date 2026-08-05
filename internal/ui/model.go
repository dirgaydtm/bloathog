package ui

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaydtm/bloathog/internal/eror"
	"github.com/dirgaydtm/bloathog/internal/monitor"
	"github.com/dirgaydtm/bloathog/internal/ring"
	"github.com/dirgaydtm/bloathog/internal/ui/components"
	"github.com/dirgaydtm/bloathog/internal/ui/theme"
	"github.com/dirgaydtm/bloathog/internal/ui/types"
)

const (
	maxGraphSamples = 120
	maxLogLines     = 5000
	focusLog        = 0
	focusProc       = 1
	focusGraph      = 2
	focusCount      = 3
)

// Model is the root app state.
type Model struct {
	projectInfo types.ProjectInfo

	// Process state
	rootPID  int32
	cmd      *exec.Cmd
	started  bool
	quitting bool
	ExitCode int

	// Interactive state
	stdin     io.WriteCloser
	inputMode bool
	textInput textinput.Model

	// Monitor state
	stats types.MonitorState

	// Dynamic buffers
	graph       *ring.Buffer[float64]
	cpuGraph    *ring.Buffer[float64]
	activeGraph int // 0 = RAM, 1 = CPU
	logs        *ring.Buffer[string]
	logDirty    bool

	// UI components
	logPanel    components.LogPanel
	procPanel   components.ProcessTreePanel
	header      components.HeaderModel
	helpModel   help.Model
	keys        types.KeyMap
	focusTarget int

	// Layout
	width  int
	height int
	procW  int // cached panel width, set in relayout()

	// Exit state
	exitReport string
	startTime  time.Time
	FatalErr   error
}

func NewModel(info types.ProjectInfo) tea.Model {
	h := help.New()
	h.Styles.ShortKey = theme.StyleAccent
	h.Styles.ShortDesc = theme.StyleMuted
	h.Styles.ShortSeparator = theme.StyleMuted

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter text to send..."
	ti.Prompt = " > "
	ti.PromptStyle = theme.StyleAccent
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	cmdStr := info.Command
	for _, a := range info.Args {
		cmdStr += " " + a
	}

	m := Model{
		projectInfo: info,
		header:      components.NewHeaderModel(cmdStr),
		helpModel:   h,
		textInput:   ti,
		keys:        types.DefaultKeyMap(),
		focusTarget: focusLog,
		logPanel:    components.NewLogPanel(80, 10),
		procPanel:   components.NewProcessTreePanel(40, 10),
		graph:       ring.New[float64](maxGraphSamples),
		cpuGraph:    ring.New[float64](maxGraphSamples),
		logs:        ring.New[string](maxLogLines),
	}

	m.logPanel.SetFocused(true)

	return m
}

func (m Model) Init() tea.Cmd {
	// Initialize header and start the background child process
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
		if m.inputMode {
			switch msg.Type {
			case tea.KeyEsc, tea.KeyEnter:
				if msg.Type == tea.KeyEnter && m.stdin != nil {
					m.logs.Push(">" + m.textInput.Value())
					m.logDirty = true
					io.WriteString(m.stdin, m.textInput.Value()+"\n")
				}
				m.inputMode = false
				m.textInput.Reset()
				m.relayout()
				return m, nil
			}

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		if m.focusTarget == focusLog && msg.Type == tea.KeyEnter {
			m.inputMode = true
			m.textInput.Focus()
			m.relayout()
			return m, textinput.Blink
		}

		// Switch graph
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
		// Process successfully started, store PID and begin ticking
		m.cmd = msg.Cmd
		m.rootPID = msg.RootPID
		m.stdin = msg.Stdin
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
		// Process new monitor stats every tick
		if !m.quitting {
			rss := msg.Stats.TotalRSS
			m.stats.CurrentRSS = rss
			if rss > m.stats.PeakRSS {
				m.stats.PeakRSS = rss
			}
			m.stats.RunningSumRSS += float64(rss)
			m.stats.SampleCount++
			m.stats.ActiveProcesses = msg.Stats.ProcessCount
			if msg.Stats.ProcessCount > m.stats.PeakProcesses {
				m.stats.PeakProcesses = msg.Stats.ProcessCount
			}
			m.graph.Push(float64(rss) / (1024 * 1024))

			cpu := msg.Stats.TotalCPU / float64(runtime.NumCPU())
			m.stats.CurrentCPU = cpu
			if cpu > m.stats.PeakCPU {
				m.stats.PeakCPU = cpu
			}
			m.stats.RunningSumCPU += cpu

			// EMA Smoothing for CPU Graph
			if cpu == 0 {
				m.stats.SmoothedCPU = 0
			} else {
				m.stats.SmoothedCPU = (cpu * 0.4) + (m.stats.SmoothedCPU * 0.6)
			}
			m.cpuGraph.Push(m.stats.SmoothedCPU)

			m.procPanel.UpdateNodes(msg.Stats.Nodes)
			cmds = append(cmds, monitor.TickCmd(m.rootPID))
		}

	case logBatchMsg:
		// Render batched lines to log panel
		for _, logMsg := range msg.msgs {
			m.logs.Push(logMsg.Line)
		}
		m.logDirty = true
		cmds = append(cmds, drainCmd(msg.ch))

	case monitor.ChildExitMsg:
		// Handle exit code and start teardown
		m.ExitCode = msg.ExitCode
		if m.ExitCode != 0 && m.graph.Len() == 0 {
			rawLogs := ring.Join(m.logs, "\n")
			m.FatalErr = &eror.Error{Msg: fmt.Sprintf("command '%s' failed to start properly\n%s", m.projectInfo.Command, rawLogs)}
		}
		return m.finalizeQuit()

	case monitor.ErrorMsg:
		m.FatalErr = msg.Err
		return m.finalizeQuit()

	case spinner.TickMsg, stopwatch.TickMsg, stopwatch.StartStopMsg:
		var cmd tea.Cmd
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.logDirty {
		m.logPanel.UpdateContent(m.logs)
		m.logDirty = false
	}

	if m.inputMode {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quitting {
		return m.exitReport
	}
	return m.renderLayout()
}
