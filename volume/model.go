package volume

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hichamlamine.github.io/sonik/styles"
)

type Model struct {
	volume         float64
	progressLength uint
}

func New(volume float64) Model {
	return Model{volume: volume, progressLength: 15}
}

type VolumeUpdateMsg float64

func UpdateVolume(v float64) tea.Cmd {
	return func() tea.Msg {
		return VolumeUpdateMsg(v)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// handle < and > keybindings
	switch msg := msg.(type) {
	case VolumeUpdateMsg:
		m.volume = float64(msg)
	}
	return m, nil
}

const (
	lightLine = "─"
	heavyLine = "━"
)

func (m Model) renderProgress() string {
	var s string
	filledPortion := strings.Repeat(heavyLine, int(float64(m.progressLength)*m.volume))
	filledPortion = styles.BarFilled.Render(filledPortion)

	emptyPortion := strings.Repeat(lightLine, int(m.progressLength)-int(float64(m.progressLength)*m.volume))
	emptyPortion = styles.BarEmpty.Render(emptyPortion)

	s = lipgloss.JoinHorizontal(lipgloss.Center, filledPortion, emptyPortion)
	return s
}

func (m Model) View() string {
	var s, volume, bar string

	if m.volume == 0 {
		volume = "Muted"
	} else {
		volume = fmt.Sprintf("%d%s", int(m.volume*100), "%")
	}

	bar = m.renderProgress()

	s = lipgloss.JoinHorizontal(lipgloss.Center, volume, " ", bar)

	return s
}
