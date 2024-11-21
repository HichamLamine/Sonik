package styles

import "github.com/charmbracelet/lipgloss"

var (
	GlobalStyle = lipgloss.NewStyle()
	ListStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			// Padding(1, 2).
			Width(100)
)

type AppStyles struct {
	ListStyle lipgloss.Style
}

func NewAppStyles() AppStyles {
	return AppStyles{ListStyle: ListStyle}
}
