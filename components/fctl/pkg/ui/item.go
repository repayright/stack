package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	title, desc string
}

func (i Item) GetTitle() string       { return i.title }
func (i Item) GetDescription() string { return i.desc }

// This is needed to implement list.Item interface
func (i Item) FilterValue() string { return i.title }

func NewItem(title, desc string) *Item {
	return &Item{
		title: title,
		desc:  desc,
	}
}

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 3 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(*Item)

	if !ok {
		return
	}

	str := fmt.Sprintf("%s\n%s\n", i.GetTitle(), i.GetDescription())

	fn := ItemStyle.Render

	fmt.Fprint(w, fn(str))
}
