package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: This should extend list.Model from github.com/charmbracelet/bubbles/list
type ListModel struct {
	list list.Model
}

func NewListModel(items []list.Item, delegate list.ItemDelegate, width int, height int) *ListModel {
	l := list.New(items, delegate, width, height)

	l.SetShowTitle(true)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return &ListModel{
		list: l,
	}
}

func NewDefaultListModel(items []list.Item) *ListModel {
	return NewListModel(items, ItemDelegate{}, ViewWidth, ViewHeight).WithMaxPossibleHeight().WithMaxPossibleWidth()
}

func (m ListModel) Init() tea.Cmd {
	return nil
}
func (m ListModel) View() string {
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

// The height counter depends on row count
// of the terminal
// Sooo where is the limit ?
func (m *ListModel) GetMaxPossibleHeight() int {

	// TODO: This should be dynamic
	// It should be calculed from ItemDelegate.Height()
	// 3*len(m.list.Items())
	// 3 = title + desc + breackline
	// + 2 line for the header
	return 3*len(m.list.Items()) + 2
}
func (m *ListModel) WithMaxPossibleHeight() *ListModel {
	m.list.SetHeight(m.GetMaxPossibleHeight())
	return m

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
