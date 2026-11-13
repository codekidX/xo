package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit   key.Binding
	Add    key.Binding
	Readme    key.Binding
	Setup     key.Binding
	Export    key.Binding
	Bootstrap key.Binding
	Config    key.Binding
	Delete    key.Binding
	CD        key.Binding
}

var Keys = keyMap{
	Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Enter:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:      key.NewBinding(key.WithKeys("esc", "b"), key.WithHelp("esc/b", "back")),
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Add:       key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Readme:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "readme")),
	Setup:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "repair")),
	Export:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "export")),
	Bootstrap: key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "bootstrap")),
	Config:    key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "config")),
	Delete:    key.NewBinding(key.WithKeys("backspace", "delete"), key.WithHelp("backspace", "delete")),
	CD:        key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "cd")),
}
