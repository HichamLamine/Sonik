package styles

import "github.com/charmbracelet/lipgloss"

var (
	GlobalStyle = lipgloss.NewStyle()
	ListStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder())
			// Padding(1, 2).
)
