package list

import (
	"fmt"
	"io"
	"strings"

	blist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type HorizontalItem struct {
	title, desc string
}

func NewHorizontalItem(title, desc string) *HorizontalItem {
	return &HorizontalItem{
		title: title,
		desc:  desc,
	}
}
func (i HorizontalItem) GetTitle() string       { return i.title }
func (i HorizontalItem) GetDescription() string { return i.desc }

// This is needed to implement list.Item interface
func (i HorizontalItem) FilterValue() string { return i.title }

func (i HorizontalItem) GetWidth() int {
	return len(i.title) + len(i.desc)

}

func (i *HorizontalItem) GetHeight() int {
	return 1
}

type HorizontalItemDelegate struct{}

func NewHorizontalItemDelegate() *HorizontalItemDelegate {
	return &HorizontalItemDelegate{}
}

func (d HorizontalItemDelegate) Height() int {
	return 1
}

func (d HorizontalItemDelegate) Spacing() int                             { return 0 }
func (d HorizontalItemDelegate) Update(_ tea.Msg, _ *blist.Model) tea.Cmd { return nil }
func (d HorizontalItemDelegate) Render(w io.Writer, m blist.Model, index int, item blist.Item) {
	i, ok := item.(*HorizontalItem)
	//Styles
	style := lipgloss.Color("#fff788")
	title := lipgloss.NewStyle().Foreground(style).Bold(true)

	valueStyle := lipgloss.Color("#b3cedc")
	desc := lipgloss.NewStyle().Foreground(valueStyle).Bold(false)
	if !ok {
		return
	}

	var str string
	if i.GetDescription() == "" {
		str = title.Render(i.GetTitle())
	} else {
		str = title.Render(i.GetTitle()+" ") + desc.Render(i.GetDescription())
	}

	_, err := fmt.Fprint(w, str)
	if err != nil {
		panic(err)
	}
}
func ApplyStyle(list []string, toFill string, beforeChar string) []string {
	max := modelutils.GetMaxCharPosXinCharList(list, beforeChar)

	//Styles
	style := lipgloss.Color("#fff788")
	title := lipgloss.NewStyle().Foreground(style).Bold(true)

	valueStyle := lipgloss.Color("#b3cedc")
	desc := lipgloss.NewStyle().Foreground(valueStyle).Bold(false)

	// fmt.Println("HERES")

	for i, line := range list {

		res := modelutils.FillCharBeforeChar(line, toFill, beforeChar, max)
		split := strings.Split(res, beforeChar)
		// fmt.Println(split)
		list[i] = title.Render(split[0]+beforeChar) + desc.Render(split[1])
	}

	return list
}
