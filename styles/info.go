package styles

import "github.com/charmbracelet/lipgloss"

var (
	Info      = lipgloss.NewStyle().Padding(0, 1)
	InfoTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#e79cfe"))

	BarEmpty  = lipgloss.NewStyle().Foreground(lipgloss.Color("#505050"))
	BarFilled = lipgloss.NewStyle().Foreground(lipgloss.Color("#e79cfe"))
)
