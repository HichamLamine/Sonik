package progress

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hichamlamine.github.io/sonik/styles"
)

type Model struct {
	progress time.Duration
	length   time.Duration

	percentage    float64
	progressWidth int
	width         int
}

func New(progress time.Duration, length time.Duration) Model {
	return Model{progress: progress, length: length, progressWidth: 15}
}

type ProgressTickMsg (time.Duration)
type LengthUpdateMsg (time.Duration)

func UpdateLength(length time.Duration) tea.Cmd {
	return func() tea.Msg {
		return LengthUpdateMsg(length)
	}
}

func TickProgress(progress time.Duration) tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return ProgressTickMsg(progress)
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressTickMsg:
		m.progress = time.Duration(msg)
		if m.progress != 0 {
			m.percentage = float64(m.progress) / float64(m.length)
		}
	case LengthUpdateMsg:
		m.length = time.Duration(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

const (
	lightLine = "─"
	heavyLine = "━"
)

func (m Model) RenderProgress() string {
	var s, filledPortion, emptyPortion string
	// fmt.Println(int(m.progress))
	filledPortion = strings.Repeat(heavyLine, int(m.percentage*float64(m.progressWidth)))
	filledPortion = styles.BarFilled.Render(filledPortion)

	emptyPortion = strings.Repeat(lightLine, m.progressWidth-int(m.percentage*float64(m.progressWidth)))
	emptyPortion = styles.BarEmpty.Render(emptyPortion)

	s = lipgloss.JoinHorizontal(lipgloss.Center, filledPortion, emptyPortion)
	return s
}

func (m Model) View() string {
	var s, progress, bar, length string
	progress = fmt.Sprintf("%v", m.progress.Round(time.Second))
	length = fmt.Sprintf("%v", m.length.Round(time.Second))

	if m.width != 0 {
		m.progressWidth = m.width - styles.Info.GetHorizontalFrameSize() - lipgloss.Width(progress) - lipgloss.Width(length) - 2
	}

	bar = m.RenderProgress()

	s = lipgloss.JoinHorizontal(lipgloss.Center, progress, " ", bar, " ", length)

	return s
}
