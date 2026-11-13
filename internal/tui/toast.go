package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ShowToastMsg struct {
	Message string
	ID      int64
}

type ClearToastMsg struct {
	ID int64
}

func ShowToast(msg string) tea.Cmd {
	id := time.Now().UnixNano()
	return tea.Batch(
		func() tea.Msg { return ShowToastMsg{Message: msg, ID: id} },
		tea.Tick(2*time.Second, func(time.Time) tea.Msg { return ClearToastMsg{ID: id} }),
	)
}
