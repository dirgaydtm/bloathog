package types

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for bloathog.
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	SwitchGraph key.Binding
	Tab         key.Binding
	Enter       key.Binding
	Esc         key.Binding
	Quit        key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "left"),
			key.WithHelp(" ↑/← ", "scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "right"),
			key.WithHelp(" ↓/→ ", "scroll down"),
		),
		SwitchGraph: key.NewBinding(
			key.WithKeys("h", "l", "[", "]"),
			key.WithHelp(" h/l ", "switch graph"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(" tab ", "switch panel"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(" enter ", "input mode"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(" esc ", "exit"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "Q", "ctrl+C"),
			key.WithHelp(" q ", "quit"),
		),
	}
}

// --- Dynamic KeyMap Wrappers ---

// ProcKeyMap is used when the Proc panel is focused
type ProcKeyMap struct { KeyMap }

func (k ProcKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Tab, k.Quit}
}
func (k ProcKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Tab, k.Quit},
	}
}

// LogKeyMap is used when the Log panel is focused
type LogKeyMap struct {
	KeyMap
	InputMode bool
}

func (k LogKeyMap) ShortHelp() []key.Binding {
	if k.InputMode {
		return []key.Binding{k.Esc}
	}
	return []key.Binding{k.Up, k.Down, k.Tab, k.Enter, k.Quit}
}

func (k LogKeyMap) FullHelp() [][]key.Binding {
	if k.InputMode {
		return [][]key.Binding{{k.Esc}}
	}
	return [][]key.Binding{
		{k.Up, k.Down, k.Tab, k.Enter, k.Quit},
	}
}

// GraphKeyMap is used when the Graph panel is focused
type GraphKeyMap struct { KeyMap }

func (k GraphKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SwitchGraph, k.Tab, k.Quit}
}
func (k GraphKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SwitchGraph, k.Tab, k.Quit},
	}
}
