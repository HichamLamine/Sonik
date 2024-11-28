package lyrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/progress"
	"hichamlamine.github.io/sonik/styles"
)

const url = "https://lrclib.net/api/"

type Model struct {
	playingSong           *player.Song
	spinner               spinner.Model
	viewport              viewport.Model
	lines                 []string
	timestamps            []time.Duration
	currentIndex          int
	progress              time.Duration
	status                int
	width, height         int
	termWidth, termHeight int
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func NewModel() Model {
	sp := spinner.Points
	viewport := viewport.Model{
		Width:  50,
		Height: 20,
	}
	return Model{spinner: spinner.Model{Spinner: sp}, viewport: viewport, progress: time.Duration(0), width: 40}
}

func (m Model) GetWidth() int { return m.width }

func (m *Model) SetPlayingSong(song *player.Song) {
	m.playingSong = song
}

const (
	Idle = iota
	Requesting
	Answered
)

type LyricsResponse struct {
	Title        string `json:"trackName"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

type LyricsOkMsg struct {
	timestamps []time.Duration
	lines      []string
}

type LyricsErrMsg struct {
	desc string
	err  error
}

func parseTimestamps(tsStrs []string) ([]time.Duration, error) {
	var ts []time.Duration
	var tmpStrs []string
	for _, tsStr := range tsStrs {
		tmpstr := tsStr[1 : len(tsStr)-1]
		tmpstr = fmt.Sprintf("%s0", tmpstr)
		tmpStrs = append(tmpStrs, tmpstr)
	}
	var minStrs []string
	var secStrs []string
	var msStrs []string
	var timestampsStrs []string
	for i, tmpStr := range tmpStrs {
		minStrs = append(minStrs, fmt.Sprintf("%sm", tmpStr[0:2]))
		secStrs = append(secStrs, fmt.Sprintf("%ss", tmpStr[3:5]))
		msStrs = append(msStrs, fmt.Sprintf("%sms", tmpStr[6:]))
		timestampsStrs = append(timestampsStrs, fmt.Sprintf("%s%s%s", minStrs[i], secStrs[i], msStrs[i]))

		parsedTs, err := time.ParseDuration(timestampsStrs[i])
		if err != nil {
			return nil, err
		}
		ts = append(ts, parsedTs)
	}

	return ts, nil
}

func loadLyrics(song *player.Song) tea.Cmd {
	return func() tea.Msg {
		client := http.Client{}
		artist := strings.ReplaceAll(strings.TrimSpace(song.Artist), " ", "+")
		title := strings.ReplaceAll(strings.TrimSpace(song.Title), " ", "+")
		// title = strings.Split(title, "-")[0]
		// title = strings.Trim(title, "+")
		// log.Fatal(title)
		parametered_url := fmt.Sprintf("%sget?artist_name=%s&track_name=%s", url, artist, title)
		res, err := client.Get(parametered_url)
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return LyricsErrMsg{desc: "Could not get the lyrics.", err: fmt.Errorf(res.Status)}
		}

		response := LyricsResponse{}
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return LyricsErrMsg{desc: "failed to decode response", err: err}
		}

		rawLines := strings.Split(response.SyncedLyrics, "\n")
		var timestamps []string
		var lines []string

		for _, line := range rawLines {
			splitLine := strings.Split(line, " ")
			timestamps = append(timestamps, splitLine[0])
			lines = append(lines, strings.Join(splitLine[1:], " "))
		}

		var ts []time.Duration
		ts, err = parseTimestamps(timestamps)
		if err != nil {
			return LyricsErrMsg{desc: "failed to parse timestamps", err: err}
		}

		return LyricsOkMsg{
			timestamps: ts,
			lines:      lines,
		}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.status = Requesting
			return m, loadLyrics(m.playingSong)
		}
	case LyricsOkMsg:
		m.lines = msg.lines
		m.timestamps = msg.timestamps
		m.currentIndex = 0
		lines := slices.Clone(m.lines)
		for i, line := range lines {
			lines[i] = styles.InactiveLyricsStyle.Render(line)
		}
		longestLine := slices.MaxFunc(m.lines, func(l1, l2 string) int { return len(l1) - len(l2) })
		longLineWidth := len(longestLine) + 2
		m.viewport.Width = longLineWidth
        m.viewport.Height = m.height
		m.width = longLineWidth + styles.LyricsStyle.GetHorizontalFrameSize()
		styles.LyricsStyle = styles.LyricsStyle.Width(m.width)

		m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Center, lines...))
		m.status = Answered
		return m, nil

	case LyricsErrMsg:
		m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, msg.desc, msg.err.Error()))
		m.lines = nil
		m.timestamps = nil
		m.currentIndex = 0
		m.status = Answered
		return m, nil

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		m.height = m.termHeight - 4
		styles.LyricsStyle = styles.LyricsStyle.Width(m.width)

	case progress.ProgressTickMsg:
		m.progress = time.Duration(msg)
		if m.timestamps != nil {
			if m.progress >= m.timestamps[m.currentIndex] {

				lines := slices.Clone(m.lines)
				for i, line := range lines {
					if m.currentIndex == i {
						lines[i] = styles.ActiveLyricsStyle.Render(line)
					} else {
						lines[i] = styles.InactiveLyricsStyle.Render(line)
					}
				}

				m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Center, lines...))
				if m.currentIndex > m.height/2 {
					m.viewport.LineDown(1)
				}
				if m.currentIndex < len(lines)-1 {
					m.currentIndex++
				}
			}
		}
	}

	// compare player position and timestamp

	newSpinnerModel, spinnerCmd := m.spinner.Update(msg)
	m.spinner = newSpinnerModel
	newViewportModel, viewportCmd := m.viewport.Update(msg)
	m.viewport = newViewportModel
	cmds = append(cmds, spinnerCmd, viewportCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string

	if m.status == Answered {
		s = styles.LyricsViewportStyle.Render(m.viewport.View())
	} else if m.status == Requesting {
		s = lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " Loading lyrics")
	} else if m.status == Idle {
		s = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "No song is playing")
	}
	return s
}
