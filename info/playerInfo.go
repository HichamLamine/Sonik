package info

import (
	tea "github.com/charmbracelet/bubbletea"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/progress"
)

type Model struct {
	progress progress.Model

	// player Player
	currentSong player.Song
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			// m.player.togglePause()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return ""
}
