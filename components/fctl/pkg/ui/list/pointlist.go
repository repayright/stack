package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type PointList struct {
	list list.Model
}

func NewPointList(items []list.Item, delegate list.ItemDelegate, width int, height int) *PointList {
	l := list.New(items, delegate, width, height)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return &PointList{
		list: l,
	}
}
func (pl *PointList) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}
func (pl *PointList) Init() tea.Cmd {
	return nil
}
func (pl *PointList) Update(msg tea.Msg) (*PointList, tea.Cmd) {
	var cmd tea.Cmd
	pl.list, cmd = pl.list.Update(msg)
	return pl, cmd
}

func (pl *PointList) GetMaxPossibleWidth() int {
	max := 0
	for _, i := range pl.list.Items() {

		item, ok := i.(HorizontalItem)
		if !ok {
			return 0
		}
		if item.GetWidth() >= max {
			max = item.GetWidth()
		}
	}
	return max
}

func (pl *PointList) View() string {
	return pl.list.View()
}
