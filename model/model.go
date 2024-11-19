package model

import (
	// "fmt"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	"hichamlamine.github.io/sonik/player"
)

type SongsDataMsg []player.Song
type ReadErrorMsg struct{ err error }
type MetadataErrorMsg struct {
	song string
	err  error
}

type item struct {
	title, artist string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.artist }
func (i item) FilterValue() string { return i.title }

func LoadSongs() tea.Msg {
	homedir := os.Getenv("HOME")
	musicdir := filepath.Join(homedir, "Music")

	songsdir, err := os.ReadDir(musicdir)

	if err != nil {
		return ReadErrorMsg{err}
	}

	var songsData []player.Song
	for _, file := range songsdir {
		fileExtension := file.Name()[len(file.Name())-3:]
		if !file.IsDir() && fileExtension != "zip" {
			fileDir := filepath.Join(musicdir, file.Name())
			f, err := os.Open(fileDir)
			if err != nil {
				return ReadErrorMsg{err}
			}
			m, err := tag.ReadFrom(f)
			if err != nil {
				return MetadataErrorMsg{fileDir, err}
			}

			if m.FileType() == "MP3" {
				songsData = append(songsData, player.Song{
					Path:   fileDir,
					Title:  m.Title(),
					Artist: m.Artist(),
					Lyrics: m.Lyrics(),
				})
			}

		}
	}
	return SongsDataMsg(songsData)
}

type Model struct {
	Songs []player.Song
	List  list.Model
}

func (m Model) Init() tea.Cmd {
	return LoadSongs
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case SongsDataMsg:
		m.Songs = msg
		return m, nil
	case ReadErrorMsg:
		fmt.Println("ReadErrorMsg received:", msg.err)
		return m, tea.Quit
	case MetadataErrorMsg:
		fmt.Printf("%s: MetadataErrorMsg received:%s\n", msg.song, msg.err)
		return m, nil
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width, msg.Height)
	default:
		fmt.Println("Unknown message received:", msg)
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var items []list.Item
	for _, song := range m.Songs {
		items = append(items, item{
			title:  song.Title,
			artist: song.Artist,
		})
	}
	m.List = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m.List.View()
}
