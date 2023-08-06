package list

import (
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui/helpers"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

type PointList struct {
	list []*HorizontalItem
}

func NewPointList(list ...*HorizontalItem) *PointList {
	return &PointList{
		list: list,
	}
}
func (pl *PointList) GetListKeyMapHandler() *config.KeyMapHandler {
	return nil
}

// SortDir sort the list by asc or desc
// true = asc
// false = desc
func (pl *PointList) SortDir(dir bool) {
	sort.Slice(pl.list, func(i int, j int) bool {
		if dir {
			return pl.list[i].title < pl.list[j].title
		}
		return pl.list[i].title > pl.list[j].title
	})

}

func (pl *PointList) Init() tea.Cmd {
	return nil
}
func (pl *PointList) Update(msg tea.Msg) (*PointList, tea.Cmd) {
	return pl, nil
}

func (pl *PointList) GetMaxPossibleWidth() int {
	max := 0
	for _, item := range pl.list {
		if item.GetWidth() >= max {
			max = item.GetWidth()
		}
	}
	return max
}

func (pl *PointList) View() string {
	var section = make([]string, 0)
	style := theme.YellowColor
	title := lipgloss.NewStyle().Foreground(style).Bold(true)

	valueStyle := lipgloss.Color("#b3cedc")
	desc := lipgloss.NewStyle().Foreground(valueStyle).Bold(false)

	Log := helpers.NewLogger("POINTLIST")
	Log.Log("PointList", strconv.Itoa(len(pl.list)))

	for _, item := range pl.list {
		Log.Log("PointList ITEM", strconv.FormatBool(item == nil))

		str := title.Render(item.title+" ") + desc.Render(item.desc)
		section = append(section, str)
	}
	return strings.Join(section, "\n")
}
