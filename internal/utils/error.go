package utils

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4C4C")).Bold(true)

// Fatal prints a beautifully formatted error message and exits the program cleanly.
func Fatal(msg any) {
	fmt.Println(errorStyle.Render(fmt.Sprintf("Error: %v", msg)))
	os.Exit(1)
}
