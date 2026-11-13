package tui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor   = lipgloss.Color("#04B575")
	secondaryColor = lipgloss.Color("#AEAEAE")
	errorColor     = lipgloss.Color("#FF4C4C")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 1).
			MarginBottom(1)

	formLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD")).
			Width(12)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			MarginTop(1)

	appBoxStyle = lipgloss.NewStyle().
			Margin(1, 2)

	toastStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#555555")).
			Padding(0, 1).
			MarginLeft(2)
)
