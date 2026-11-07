package tui

import (
	"os"
	"xo/internal/store"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectSelectedMsg struct{ Name string }
type OpenAddProjectFormMsg struct{}
type SetupWorkspaceMsg struct{}
type ExportWorkspaceMsg struct{}
type OpenBootstrapFormMsg struct{}
type OpenSetupListMsg struct{}
type CDProjectMsg struct{ Path string }

type projectItem struct {
	name    string
	path    string
	missing bool
}

func (i projectItem) Title() string { 
	if i.missing {
		return lipgloss.NewStyle().Foreground(errorColor).Render(i.name)
	}
	return i.name 
}
func (i projectItem) Description() string { 
	if i.missing {
		return lipgloss.NewStyle().Foreground(errorColor).Render(i.path + " (missing)")
	}
	return i.path 
}
func (i projectItem) FilterValue() string { return i.name }

type projectsModel struct {
	list list.Model
}

func newProjectsModel(s *store.Store) projectsModel {
	projs := s.GetProjects()
	items := make([]list.Item, 0, len(projs))
	for name, p := range projs {
		_, err := os.Stat(p)
		items = append(items, projectItem{name: name, path: p, missing: os.IsNotExist(err)})
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{Keys.Add, Keys.Readme, Keys.CD}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{Keys.Add, Keys.Readme, Keys.Setup, Keys.Export, Keys.Bootstrap, Keys.Config, Keys.CD}
	}

	return projectsModel{list: l}
}

func (m *projectsModel) resize(w, h int) {
	// subtract appBox margins (2 vertically) + explicit title height (3 lines)
	m.list.SetSize(w-4, h-5)
}

func (m projectsModel) Update(msg tea.Msg) (projectsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, Keys.Add):
				return m, func() tea.Msg { return OpenAddProjectFormMsg{} }
			case key.Matches(msg, Keys.CD):
				if i, ok := m.list.SelectedItem().(projectItem); ok {
					if !i.missing {
						return m, func() tea.Msg { return CDProjectMsg{Path: i.path} }
					}
				}
			case key.Matches(msg, Keys.Setup):
				return m, func() tea.Msg { return SetupWorkspaceMsg{} }
			case key.Matches(msg, Keys.Export):
				return m, func() tea.Msg { return ExportWorkspaceMsg{} }
			case key.Matches(msg, Keys.Bootstrap):
				return m, func() tea.Msg { return OpenBootstrapFormMsg{} }
			case key.Matches(msg, Keys.Config):
				return m, func() tea.Msg { return OpenSetupListMsg{} }
			case key.Matches(msg, Keys.Readme):
				if i, ok := m.list.SelectedItem().(projectItem); ok {
					return m, func() tea.Msg { return OpenReadmeMsg{Name: i.name, Path: i.path} }
				}
			case key.Matches(msg, Keys.Enter):
				if i, ok := m.list.SelectedItem().(projectItem); ok {
					return m, func() tea.Msg { return ProjectSelectedMsg{Name: i.name} }
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m projectsModel) View() string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("XO Projects"),
		m.list.View(),
	)
	return appBoxStyle.Render(content)
}
