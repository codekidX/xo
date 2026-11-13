package tui

import (
	"xo/internal/types"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommandSelectedMsg struct {
	CmdStr  string
	WorkDir string
	Path    string
}

type OpenAddFormMsg struct{}

type commandItem struct {
	cmd types.Command
}

func (i commandItem) Title() string       { return i.cmd.Name }
func (i commandItem) Description() string { 
	if i.cmd.Help != "" {
		return i.cmd.Help + " [" + i.cmd.CmdStr + "]" 
	}
	return "[" + i.cmd.CmdStr + "]"
}
func (i commandItem) FilterValue() string { return i.cmd.Name }

type commandsModel struct {
	projectName string
	projectPath string
	list        list.Model
}

func newCommandsModel(name, path string, cmdMap types.CommandMap) commandsModel {
	items := make([]list.Item, 0, len(cmdMap))
	for _, c := range cmdMap {
		items = append(items, commandItem{cmd: c})
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{Keys.Add, Keys.Back}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{Keys.Add, Keys.Back}
	}

	return commandsModel{projectName: name, projectPath: path, list: l}
}

func (m *commandsModel) resize(w, h int) {
	// subtract margins and explicit title height (3 lines)
	m.list.SetSize(w-4, h-5)
}

func (m commandsModel) Update(msg tea.Msg) (commandsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, Keys.Add):
				return m, func() tea.Msg { return OpenAddFormMsg{} }
			case key.Matches(msg, Keys.Enter):
				if i, ok := m.list.SelectedItem().(commandItem); ok {
					return m, func() tea.Msg {
						return CommandSelectedMsg{
							CmdStr:  i.cmd.CmdStr,
							WorkDir: i.cmd.WorkDir,
							Path:    m.projectPath,
						}
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m commandsModel) View() string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Commands"),
		m.list.View(),
	)
	return appBoxStyle.Render(content)
}
