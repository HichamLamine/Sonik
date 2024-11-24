package info

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/progress"
	"hichamlamine.github.io/sonik/styles"
	"hichamlamine.github.io/sonik/volume"
)

type Model struct {
	progress progress.Model
	volume   volume.Model

	player *player.State

	width int
}

func (m Model) Init() tea.Cmd {
	return progress.TickProgress(m.player.GetPosition())
}

// p pointer
func New(p *player.State) Model {
	// init progress
	// if p == nil {
	v := volume.New(p.GetVolume())
	progress := progress.New(p.GetPosition(), p.GetLen())
	return Model{player: p, volume: v, width: 0, progress: progress}
}

// func(m Model) UpdateProgress() tea.Cmd {
//     return func() tea.Msg {
//
//     }
// }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case " ":
		// 	m.player.TogglePause()
		case "enter":
			return m, progress.UpdateLength(m.player.GetLen())
		case "<", ">":
			return m, volume.UpdateVolume(m.player.GetVolume())
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case progress.ProgressTickMsg:
		cmds = append(cmds, progress.TickProgress(m.player.GetPosition()))
	}

	volume, cmd := m.volume.Update(msg)
	m.volume = volume
	cmds = append(cmds, cmd)

	progress, cmd := m.progress.Update(msg)
	m.progress = progress
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

const (
	PlayIcon  = "󰐊"
	PauseIcon = "󰏤"
)

func (m Model) View() string {
	var s string

	var title, icon, volume string
	if m.player.PlayingSong == nil {
		title = "No song is playing"
	} else {
		title = m.player.PlayingSong.Title
	}
	if m.player.IsPaused() {
		icon = PauseIcon
	} else {
		icon = PlayIcon
	}

	songName := lipgloss.JoinHorizontal(lipgloss.Center, icon, "  ", title)
	songName = styles.InfoTitle.Render(songName)

	volume = m.volume.View()

	volume = lipgloss.PlaceHorizontal(m.width-styles.Info.GetHorizontalFrameSize()-lipgloss.Width(songName), lipgloss.Right, volume)
	// fmt.Println(styles.Info.GetHorizontalFrameSize())

	progress := m.progress.View()

	s = progress

	s = lipgloss.JoinVertical(lipgloss.Center, s, lipgloss.JoinHorizontal(lipgloss.Center, songName, volume))

	s = styles.Info.Render(s)

	return s
}
