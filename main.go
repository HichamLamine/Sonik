package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/info"
	"hichamlamine.github.io/sonik/lyrics"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/styles"
	"hichamlamine.github.io/sonik/utils"
)

type Model struct {
	// state sessionState
	// pages []page.Model
	list   list.Model
	info   info.Model
	lyrics lyrics.Model

	songs []player.Song
	p     *player.State

	focusedModel int

	lyricsWidth   int
	width, height int
}

const (
	List = iota
	Lyrics
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func newModel() Model {
	var items []list.Item
	songs := utils.LoadSongs()
	slices.SortFunc(songs, func(s1, s2 player.Song) int {
		return strings.Compare(strings.ToLower(s1.Title), strings.ToLower(s2.Title))
	})
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

	lyricsModel := lyrics.NewModel()

	return Model{list: listModel, info: infoModel, songs: songs, p: &p, lyrics: lyricsModel, lyricsWidth: lyricsModel.GetWidth()}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.info.Init(), m.lyrics.Init())
}

func (m *Model) ToggleFocus() {
	if m.focusedModel == List {
		m.focusedModel = Lyrics
		styles.LyricsStyle = styles.LyricsStyle.BorderForeground(lipgloss.Color("#e79cfe"))
		styles.ListStyle = styles.ListStyle.BorderForeground(lipgloss.Color("#fff"))
	} else {
		m.focusedModel = List
		styles.ListStyle = styles.ListStyle.BorderForeground(lipgloss.Color("#e79cfe"))
		styles.LyricsStyle = styles.LyricsStyle.BorderForeground(lipgloss.Color("#fff"))
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// fmt.Println(m.list.SelectedItem())
			if m.list.FilterState() == list.FilterApplied {
				for _, song := range m.songs {
					if fmt.Sprintf("{%s}", strings.Join([]string{song.Title, song.Artist}, " ")) == fmt.Sprintln(m.list.SelectedItem()) {
						m.p.Play(&song)
					}
				}
			}
			if m.list.FilterState() == list.Unfiltered {
				song := m.songs[m.list.Index()]
				m.p.Play(&song)
				m.lyrics.SetPlayingSong(&song)
			}
		case " ":
			m.p.TogglePause()
		case "<":
			m.p.DecreaseVolume()
		case ">":
			m.p.IncreaseVolume()

		case "ctrl+l", "ctrl+right":
			m.ToggleFocus()

		case "ctrl+h", "ctrl+left":
			m.ToggleFocus()
		}

	// case lyrics.LyricsOkMsg:

	case lyrics.LyricsWidthMsg:
		infoHeight := lipgloss.Height(m.info.View())
		m.lyricsWidth = int(msg)
		m.list.SetSize(m.width-m.lyricsWidth-styles.ListStyle.GetHorizontalFrameSize(), m.height-2-infoHeight)
		styles.ListStyle = styles.ListStyle.Width(m.width - m.lyricsWidth - 4)

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		infoHeight := lipgloss.Height(m.info.View())
		m.list.SetSize(m.width-m.lyricsWidth-2, m.height-2-infoHeight)
		styles.ListStyle = styles.ListStyle.Width(m.width - m.lyricsWidth - 4)
		// fmt.Println(m.list.Width())

		// styles.LyricsStyle = styles.LyricsStyle.Height(m.height - infoHeight - 2)
		// m.list.SetSize(w-2, h-2-infoHeight)
		// styles.ListStyle = styles.ListStyle.Width(m.list.Width()).Height(m.list.Height() - infoHeight)
	}
	newListModel, listCmd := m.list.Update(msg)
	newInfoModel, infoCmd := m.info.Update(msg)
	newLyricsModel, lyricsCmd := m.lyrics.Update(msg)
	cmds = append(cmds, listCmd, infoCmd, lyricsCmd)
	m.list = newListModel
	m.info = newInfoModel
	m.lyrics = newLyricsModel
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	styledList := styles.ListStyle.Render(m.list.View())
	styledLyrics := styles.LyricsStyle.Render(m.lyrics.View())
	listLyricsSection := lipgloss.JoinHorizontal(lipgloss.Center, styledList, styledLyrics)
	playerInfo := m.info.View()
	// progress := m.progress.view()
	s = lipgloss.JoinVertical(lipgloss.Left, listLyricsSection, playerInfo)
	return s
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
