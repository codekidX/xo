package tui

import (
	"xo/internal/store"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type OpenAddSetupFormMsg struct{}
type DeleteSetupCmdMsg struct {
	OSName string
	CmdStr string
}
type RunSetupCmdMsg struct {
	CmdStr string
}

type setupItem struct {
	osName string
	cmdStr string
}

func (i setupItem) Title() string       { return "[" + i.osName + "]" }
func (i setupItem) Description() string { return i.cmdStr }
func (i setupItem) FilterValue() string { return i.osName + " " + i.cmdStr }

type setupListModel struct {
	list list.Model
}

func newSetupListModel(s *store.Store) setupListModel {
	setupCmds := s.GetSetupCommands()
	var items []list.Item
	for osName, cmds := range setupCmds {
		for _, c := range cmds {
			items = append(items, setupItem{osName: osName, cmdStr: c})
		}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Global Setup Commands"
	l.SetShowStatusBar(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{Keys.Add, Keys.Delete, Keys.Enter}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{Keys.Add, Keys.Delete, Keys.Enter, Keys.Back}
	}

	return setupListModel{list: l}
}

func (m *setupListModel) resize(w, h int) {
	m.list.SetSize(w-4, h-2)
}

func (m setupListModel) Update(msg tea.Msg) (setupListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, Keys.Add):
				return m, func() tea.Msg { return OpenAddSetupFormMsg{} }
			case key.Matches(msg, Keys.Delete):
				if i, ok := m.list.SelectedItem().(setupItem); ok {
					return m, func() tea.Msg { return DeleteSetupCmdMsg{OSName: i.osName, CmdStr: i.cmdStr} }
				}
			case key.Matches(msg, Keys.Enter):
				if i, ok := m.list.SelectedItem().(setupItem); ok {
					return m, func() tea.Msg { return RunSetupCmdMsg{CmdStr: i.cmdStr} }
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m setupListModel) View() string {
	return appBoxStyle.Render(m.list.View())
}
