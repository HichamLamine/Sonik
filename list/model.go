package list

import (
	"fmt"

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
	return Model{
		Items:    list,
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

type SelectItemMsg string

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			previousItemIndex := m.Selected.Index
			if previousItemIndex == 0 {
				previousItemIndex = len(m.Items)
			}
			m.Selected.Index = previousItemIndex - 1
			m.Selected.Item = m.Items[previousItemIndex-1]
			return m, nil
		case "j":
			previousItemIndex := m.Selected.Index
			if previousItemIndex == len(m.Items)-1 {
				previousItemIndex = -1
			}
			m.Selected.Index = previousItemIndex + 1
			m.Selected.Item = m.Items[previousItemIndex+1]
			return m, nil

		case "enter":
			m.SelectedFunc(m.Selected)
		}
	}
	return m, nil
}

func (m Model) View() string {
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

	return s
}
