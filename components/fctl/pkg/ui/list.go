package ui

import (
	"errors"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	fctl "github.com/formancehq/fctl/pkg"
)

// TODO: This should extend list.Model from github.com/charmbracelet/bubbles/list
type ListModel struct {
	list list.Model
	help bool
}

func NewListModel(items []list.Item, delegate list.ItemDelegate, width int, height int, help bool) *ListModel {
	l := list.New(items, delegate, width, height)

	l.SetShowTitle(true)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return &ListModel{
		list: l,
		help: help,
	}
}

// ViewWidth, ViewHeight
// Default width and height
// Should be dynamic and scale with terminale view
func NewDefaultListModel(items []list.Item, help bool) (*ListModel, error) {
	if len(items) == 0 {
		return nil, errors.New("ITEMS_EMPTY")
	}

	firstItem, ok := items[0].(*Item)
	if !ok {
		return nil, errors.New("FIRST_ITEMS_NOT_ITEM")
	}

	m := NewListModel(items, NewItemDelegate(firstItem.GetHeight()), ViewWidth, ViewHeight, help).WithMaxPossibleWidth()

	m, err := m.WithMaxPossibleHeight()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) helpView() string {
	return HelpStyle.Render("Formance CLI: \n • ↑/↓: Navigate \n • q: Quit")
}

func (m ListModel) View() string {
	if m.help {
		return m.helpView() + "\n" + DocStyle.Render(m.list.View())
	}

	return DocStyle.Render(m.list.View())
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ListModel) WithTitle(title string) *ListModel {
	m.list.Title = title
	return m
}

// The width counter depends on the max character
// of the longest line of the terminal
// The terminal width is limitant
func (m *ListModel) GetMaxPossibleWidth() int {
	return 90
}

// header is equivalent to one line + 1 breackline
func (m *ListModel) GetHeaderHeight() int {
	if m.list.ShowTitle() {
		return 2
	}
	return 0
}

func (m *ListModel) GetFooterHeight() int {
	return 0
}

// Each item has X lines defined with ItemDelegate.Height()
// Each item has 1 breackline
// It should be calculed from ItemDelegate.Height()
func (m *ListModel) GetBodyHeight() (int, error) {
	itemListLength, err := fctl.Reduce[list.Item, int](m.list.Items(), func(acc int, i list.Item, err error) (int, error) {
		item, ok := i.(*Item)
		if !ok {
			return acc, errors.New("ITEM_NOT_ITEM")
		}

		return acc + item.GetHeight(), nil
	}, 0)

	if err != nil {
		return 0, err
	}

	return itemListLength + len(m.list.Items()), nil
}

// The height counter depends on row count
// of the terminal
// res = header + body + footer
func (m *ListModel) GetMaxPossibleHeight() (int, error) {
	bodyHeight, err := m.GetBodyHeight()
	if err != nil {
		return 0, err
	}

	return m.GetHeaderHeight() + bodyHeight + m.GetFooterHeight(), nil
}

func (m *ListModel) WithMaxPossibleHeight() (*ListModel, error) {
	height, err := m.GetMaxPossibleHeight()
	if err != nil {
		return nil, err
	}

	m.list.SetHeight(height)
	return m, nil

}
func (m *ListModel) WithMaxPossibleWidth() *ListModel {
	m.list.SetWidth(m.GetMaxPossibleWidth())

	return m
}
func (m ListModel) GetHeigth() int {
	return m.list.Height()
}
func (m ListModel) GetWidth() int {
	return m.list.Width()
}
