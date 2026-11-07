package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BootstrapFormSubmitMsg struct {
	URLOrPath string
}

type setupFormModel struct {
	input        textinput.Model
	errorMessage string
	focusButton  bool
}

func newSetupFormModel() setupFormModel {
	input := textinput.New()
	input.Placeholder = "https://... or path/to/workspace.json"
	input.Focus()
	input.Prompt = ""
	return setupFormModel{input: input}
}

func (m setupFormModel) Update(msg tea.Msg) (setupFormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return FormCancelMsg{} }
		case "tab", "shift+tab", "up", "down":
			m.focusButton = !m.focusButton
			if m.focusButton {
				m.input.Blur()
			} else {
				m.input.Focus()
			}
		case "enter":
			if m.focusButton {
				val := strings.TrimSpace(m.input.Value())
				if val != "" {
					return m, func() tea.Msg { return BootstrapFormSubmitMsg{URLOrPath: val} }
				}
			} else {
				m.focusButton = true
				m.input.Blur()
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m setupFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Bootstrap Shared Workspace"))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left,
		formLabelStyle.Render("URL/Path:"),
		m.input.View(),
	))
	b.WriteString("\n\n")

	button := "[ Start ]"
	if m.focusButton {
		button = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(button)
	} else {
		button = lipgloss.NewStyle().Foreground(secondaryColor).Render(button)
	}
	b.WriteString(button)
	b.WriteString("\n\n")

	if m.errorMessage != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(errorColor).Render(m.errorMessage) + "\n\n")
	}

	b.WriteString(helpStyle.Render("tab/↑/↓ navigate • enter submit • esc cancel"))
	return appBoxStyle.Render(b.String())
}
