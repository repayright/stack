package list

import (
	"errors"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

// Model TODO: This should extend list.Model from github.com/charmbracelet/bubbles/list
type Model struct {
	list     list.Model
	itemType string // Should be enum (vertical, horizontal)
}

func NewListModel(items []list.Item, delegate list.ItemDelegate, width int, height int) *Model {
	l := list.New(items, delegate, width, height)

	// Can be done in Controller for each setup
	l.SetShowTitle(true)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return &Model{
		list: l,
	}
}

// NewDefaultListModel ViewWidth, ViewHeight
// Default width and height
// Should be dynamic and scale with terminal view
func NewDefaultListModel(items []list.Item) (*Model, error) {
	if len(items) == 0 {
		return nil, errors.New("ITEMS_EMPTY")
	}

	firstItem, ok := items[0].(*VerticalItem)
	if !ok {
		return nil, errors.New("FIRST_ITEMS_NOT_ITEM")
	}

	m := NewListModel(items, NewItemDelegate(firstItem.GetHeight()), theme.ViewWidth, theme.ViewHeight).WithMaxPossibleWidth()

	m, err := m.WithMaxPossibleHeight()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return theme.DocStyle.Render(m.list.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) WithTitle(title string) *Model {
	m.list.Title = title
	return m
}

// The width counter depends on the max character
// of the longest line of the terminal
// The terminal width is limitant
func (m *Model) GetMaxPossibleWidth() int {

	max := 0
	for _, item := range m.list.Items() {
		i, ok := item.(*VerticalItem)
		if !ok {
			return 0
		}

		if w := i.GetWidth(); w >= max {
			max = w
		}
	}

	return max
}

// header is equivalent to one line + 1 breackline
func (m *Model) GetHeaderHeight() int {
	if m.list.ShowTitle() {
		return 2
	}
	return 0
}

func (m *Model) GetFooterHeight() int {
	return 0
}

// Each item has X lines defined with ItemDelegate.Height()
// Each item has 1 breackline
// It should be calculed from ItemDelegate.Height()
func (m *Model) GetBodyHeight() (int, error) {

	sum := 0

	for _, item := range m.list.Items() {
		i, ok := item.(*VerticalItem)
		if !ok {
			return 0, errors.New("ITEM_NOT_ITEM")
		}

		sum += i.GetHeight()
	}

	return sum + len(m.list.Items()), nil
}

// The height counter depends on row count
// of the terminal
// res = header + body + footer
func (m *Model) GetMaxPossibleHeight() (int, error) {
	bodyHeight, err := m.GetBodyHeight()
	if err != nil {
		return 0, err
	}

	return m.GetHeaderHeight() + bodyHeight + m.GetFooterHeight(), nil
}

func (m *Model) WithMaxPossibleHeight() (*Model, error) {
	height, err := m.GetMaxPossibleHeight()
	if err != nil {
		return nil, err
	}

	m.list.SetHeight(height)
	return m, nil

}
func (m *Model) WithMaxPossibleWidth() *Model {
	m.list.SetWidth(m.GetMaxPossibleWidth())

	return m
}
func (m Model) GetHeigth() int {
	return m.list.Height()
}
func (m Model) GetWidth() int {
	return m.list.Width()
}
