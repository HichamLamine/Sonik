package main

import (
	"log"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hichamlamine.github.io/sonik/list"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/utils"
)

type sessionState uint

const (
	listView sessionState = iota
	progressView
)

var (
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color(241))
	focusedModelStyle = lipgloss.
				NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBackground(lipgloss.Color(69))
	modelStyle = lipgloss.
			NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
)

type model struct {
	// state sessionState
	list list.Model
}

func newModel() model {
	var items []list.Item
	songs := utils.LoadSongs()
	for _, song := range songs {
		items = append(items, list.Item{Title: song.Title, Desc: song.Artist})
	}
	listModel := list.NewModel(items)
	listModel.SetSelectedFunc(func(selectedItem list.SelectedItem) {
		p := player.Player{}
		p.Play(&songs[selectedItem.Index])
	})
	return model{list: listModel}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

func (m model) View() string {
	var s string
	s += m.list.View()
	return s
}

const (
	playIcon     string = ""
	pauseIcon    string = "󰏤"
	nextIcon     string = ""
	previousIcon string = ""
)

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
