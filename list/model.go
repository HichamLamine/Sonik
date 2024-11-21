package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Title string
	Desc  string
}

type SelectedItem struct {
	Item  Item
	Index int
}

type styles struct {
	itemTitle         lipgloss.Style
	itemDesc          lipgloss.Style
	selectedItemTitle lipgloss.Style
	selectedItemDesc  lipgloss.Style
	list              lipgloss.Style
}

type Model struct {
	Items        []Item
	Selected     SelectedItem
	SelectedFunc func(selected SelectedItem)
	viewport     viewport.Model
	Styles       styles
}

func NewModel(list []Item) Model {
	itemTitleStyle := lipgloss.NewStyle()
	selectedItemTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e79cfe")).
		Bold(true)

	itemDescStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#404040"))
	selectedItemDescStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#53375c"))

	viewport := viewport.New(20, 20)
	// viewport.HighPerformanceRendering = true
	return Model{
		Items:    list,
		viewport: viewport,
		Selected: SelectedItem{Item: list[0], Index: 0},
		Styles: styles{
			itemTitle:         itemTitleStyle,
			selectedItemTitle: selectedItemTitleStyle,
			itemDesc:          itemDescStyle,
			selectedItemDesc:  selectedItemDescStyle,
		},
	}
}

func (m *Model) SetSelectedFunc(f func(selectedItem SelectedItem)) {
	m.SelectedFunc = f
}

func (m *Model) SelectItem(index int) tea.Cmd {
	return func() tea.Msg {
		m.Items = append(m.Items, Item{Title: "a new one"})
		m.Selected = SelectedItem{
			Item:  m.Items[index],
			Index: index,
		}
		return SelectItemMsg(fmt.Sprintf("selected item %d: %s", index, m.Items[index].Title))
	}
}

func (m *Model) UpdateViewportContent() {
	var s string
	for _, item := range m.Items {
		if m.Selected.Item == item {
			s = lipgloss.JoinVertical(lipgloss.Left,
				s,
				m.Styles.selectedItemTitle.Render(item.Title),
				m.Styles.selectedItemDesc.Render(item.Desc),
				"\n",
			)
		} else {
			s = lipgloss.JoinVertical(lipgloss.Left,
				s,
				m.Styles.itemTitle.Render(item.Title),
				m.Styles.itemDesc.Render(item.Desc),
				"\n",
			)
		}
	}

	m.viewport.SetContent(s)
}

type SelectItemMsg string

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			if m.Selected.Index > 0 {
				m.Selected.Index--
			} else {
				m.Selected.Index = len(m.Items) - 1
				m.viewport.GotoBottom()
			}
			m.Selected.Item = m.Items[m.Selected.Index]
			if m.Selected.Index < m.viewport.YOffset {
				m.viewport.LineUp(4)
			}
			m.UpdateViewportContent()
			return m, nil
		case "j":
			if m.Selected.Index < len(m.Items)-1 {
				m.Selected.Index++
			} else {
				m.viewport.GotoTop()
				m.Selected.Index = 0
			}
			m.Selected.Item = m.Items[m.Selected.Index]
			if m.Selected.Index >= m.viewport.YOffset+m.viewport.Height {
				m.viewport.LineDown(4)
			}
			m.UpdateViewportContent()
			return m, nil

		case "enter":
			m.SelectedFunc(m.Selected)
		}
	case tea.WindowSizeMsg:
		m.UpdateViewportContent()
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	return m.viewport.View()
}
