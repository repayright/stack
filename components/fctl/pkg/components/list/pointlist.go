package list

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/theme"
)

type PointListOpts func(*PointList) *PointList
type PointList struct {
	style lipgloss.Style

	list []*HorizontalItem
}

func NewPointList(list []*HorizontalItem, opts ...PointListOpts) *PointList {
	pl := &PointList{
		list: list,
	}

	for _, opt := range opts {
		pl = opt(pl)
	}

	return pl
}

func WithPointListStyle(style lipgloss.Style) PointListOpts {
	return func(pl *PointList) *PointList {
		pl.style = style
		return pl
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

	for _, item := range pl.list {
		str := title.Render(item.title+" ") + desc.Render(item.desc)
		section = append(section, str)
	}
	return strings.Join(section, "\n")
}
