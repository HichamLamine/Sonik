package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/info"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/styles"
	"hichamlamine.github.io/sonik/utils"
)

type model struct {
	// state sessionState
	// pages []page.Model
	list list.Model
	info info.Model

	songs []player.Song
	p     *player.State
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func newModel() model {
	var items []list.Item
	songs := utils.LoadSongs()
	for _, song := range songs {
		items = append(items, item{title: song.Title, desc: song.Artist})
	}

	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.SetShowTitle(false)
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(false)

	p, err := player.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}

	infoModel := info.New(&p)

	return model{list: listModel, info: infoModel, songs: songs, p: &p}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.list.FilterState() != list.Filtering {
				song := m.songs[m.list.Index()]
				m.p.Play(&song)
			}
		case " ":
			m.p.TogglePause()
		case "<":
			m.p.DecreaseVolume()
		case ">":
			m.p.IncreaseVolume()
		}
	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		infoHeight := lipgloss.Height(m.info.View())
		m.list.SetSize(w-2, h-2-infoHeight)
		styles.ListStyle = styles.ListStyle.Width(m.list.Width()).Height(m.list.Height() - infoHeight)
	}
	newListModel, listCmd := m.list.Update(msg)
	newInfoModel, infoCmd := m.info.Update(msg)
	cmds = append(cmds, listCmd, infoCmd)
	m.list = newListModel
	m.info = newInfoModel
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	styledList := styles.ListStyle.Render(m.list.View())
	playerInfo := m.info.View()
	// progress := m.progress.view()
	s = lipgloss.JoinVertical(lipgloss.Left, styledList, playerInfo)
	return s
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
