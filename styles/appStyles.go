package styles

import "github.com/charmbracelet/lipgloss"

var (
	GlobalStyle = lipgloss.NewStyle()
	ListStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#e79cfe"))
			// Padding(1, 2).
)
