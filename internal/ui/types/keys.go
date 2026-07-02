package types

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for bloathog.
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	SwitchGraph key.Binding
	Tab         key.Binding
	Help        key.Binding
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
			key.WithKeys("left", "right", "h", "l"),
			key.WithHelp(" ←/→ ", "switch graph"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(" tab ", "switch panel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "Q", "ctrl+C"),
			key.WithHelp(" q ", "quit"),
		),
	}
}

// --- Dynamic KeyMap Wrappers ---

// ScrollKeyMap is used when the Log or Proc panel is focused
type ScrollKeyMap struct { KeyMap }

func (k ScrollKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Tab, k.Quit}
}
func (k ScrollKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Tab, k.Quit},
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
