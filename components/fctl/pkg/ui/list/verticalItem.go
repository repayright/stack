package list

import (
	"fmt"
	"io"
	"math"

	blist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

type VerticalItem struct {
	title, desc string
}

func NewItem(title, desc string) *VerticalItem {
	return &VerticalItem{
		title: title,
		desc:  desc,
	}
}
func (i VerticalItem) GetTitle() string       { return i.title }
func (i VerticalItem) GetDescription() string { return i.desc }

// This is needed to implement list.Item interface
func (i VerticalItem) FilterValue() string { return i.title }

func (i VerticalItem) GetWidth() int {

	return int(math.Max(float64(len(i.title)), float64(len(i.desc))))

}

func (i *VerticalItem) GetHeight() int {
	if i.desc != "" {
		return 2
	}
	return 1
}

type ItemDelegate struct {
	height int
}

func NewItemDelegate(heigth int) *ItemDelegate {
	return &ItemDelegate{
		height: heigth,
	}
}

func (d ItemDelegate) Height() int {
	return d.height
}

func (d *ItemDelegate) SetHeight(i int) {
	d.height = i
}
func (d ItemDelegate) Spacing() int                             { return 1 }
func (d ItemDelegate) Update(_ tea.Msg, _ *blist.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m blist.Model, index int, item blist.Item) {
	i, ok := item.(*VerticalItem)

	if !ok {
		return
	}

	var str string
	if i.GetDescription() == "" {
		str = fmt.Sprint(i.GetTitle())
	} else {
		str = fmt.Sprintf("%s\n%s", i.GetTitle(), i.GetDescription())
	}

	str = theme.ItemStyle.Render(str)

	_, err := fmt.Fprint(w, str)
	if err != nil {
		panic(err)
	}
}
