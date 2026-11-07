package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetupFormSubmitMsg struct {
	OSName string
	CmdStr string
}

type addSetupFormModel struct {
	inputs []textinput.Model
	focus  int
}

func newAddSetupFormModel() addSetupFormModel {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "e.g. darwin, linux"
	inputs[0].Focus()
	inputs[0].Prompt = ""

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "e.g. brew install go"
	inputs[1].Prompt = ""

	return addSetupFormModel{
		inputs: inputs,
		focus:  0,
	}
}

func (m addSetupFormModel) Update(msg tea.Msg) (addSetupFormModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return FormCancelMsg{} }

		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focus--
				if m.focus < 0 {
					m.focus = len(m.inputs)
				}
			} else {
				m.focus++
				if m.focus > len(m.inputs) {
					m.focus = 0
				}
			}

			for i := range m.inputs {
				if i == m.focus {
					cmds = append(cmds, m.inputs[i].Focus())
				} else {
					m.inputs[i].Blur()
				}
			}

		case "enter":
			if m.focus == len(m.inputs) {
				return m, func() tea.Msg {
					return SetupFormSubmitMsg{
						OSName: strings.TrimSpace(m.inputs[0].Value()),
						CmdStr: strings.TrimSpace(m.inputs[1].Value()),
					}
				}
			}

			m.focus++
			for i := range m.inputs {
				if i == m.focus {
					cmds = append(cmds, m.inputs[i].Focus())
				} else {
					m.inputs[i].Blur()
				}
			}
		}
	}

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m addSetupFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Global Setup Command"))
	b.WriteString("\n\n")

	labels := []string{"OS Name:", "Command:"}

	for i := range m.inputs {
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left,
			formLabelStyle.Render(labels[i]),
			m.inputs[i].View(),
		))
		b.WriteString("\n\n")
	}

	button := "[ Submit ]"
	if m.focus == len(m.inputs) {
		button = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(button)
	} else {
		button = lipgloss.NewStyle().Foreground(secondaryColor).Render(button)
	}
	b.WriteString(button)
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("tab/↑/↓ navigate • enter next/submit • esc cancel"))

	return appBoxStyle.Render(b.String())
}
