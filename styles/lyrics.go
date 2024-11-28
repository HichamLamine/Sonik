package styles

import "github.com/charmbracelet/lipgloss"

var (
	LyricsStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
            Width(73)
	LyricsViewportStyle = lipgloss.NewStyle()
	// Height(80)
	ActiveLyricsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true)
	InactiveLyricsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#656565"))
)
