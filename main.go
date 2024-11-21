package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	// "hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/styles"
	"hichamlamine.github.io/sonik/utils"
)

type model struct {
	// state sessionState
	// pages []page.Model
	list list.Model

	styles *styles.AppStyles

	songs []player.Song
	p     player.Player
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

	appStyles := styles.NewAppStyles()
	// listModel.SetSelectedFunc(func(selectedItem list.SelectedItem) {
	// 	p := player.Player{}
	// 	p.Play(&songs[selectedItem.Index])
	// })
	p := player.Player{}
	return model{list: listModel, songs: songs, p: p, styles: &appStyles}
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
		case "enter":
			if m.list.FilterState() != list.Filtering {
				m.p.Play(&m.songs[m.list.Index()])
			}
		}
	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.list.SetSize(w - 2, h-2)
		m.styles.ListStyle = m.styles.ListStyle.Width(m.list.Width()).Height(m.list.Height())
	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

func (m model) View() string {
	var s string
	styledList := m.styles.ListStyle.Render(m.list.View())
    // progress := m.progress.view()
	s = styledList
	return s
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
