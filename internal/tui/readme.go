package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type OpenReadmeMsg struct {
	Name string
	Path string
}

type readmeModel struct {
	projectName string
	viewport    viewport.Model
}

func newReadmeModel(projectName, content string, width, height int) readmeModel {
	// Render markdown with Glamour
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width-4),
	)
	str, _ := renderer.Render(content)

	vp := viewport.New(width-4, height-6)
	vp.SetContent(str)

	return readmeModel{
		projectName: projectName,
		viewport:    vp,
	}
}

func (m *readmeModel) resize(w, h int) {
	m.viewport.Width = w - 4
	m.viewport.Height = h - 6
}

func (m readmeModel) Update(msg tea.Msg) (readmeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" {
			return m, func() tea.Msg { return FormCancelMsg{} }
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m readmeModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("README · " + m.projectName))
	b.WriteString("\n\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓ scroll • q/esc back"))
	return appBoxStyle.Render(b.String())
}
