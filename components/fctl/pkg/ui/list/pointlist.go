package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"strings"
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
	style := lipgloss.Color("#fff788")
	title := lipgloss.NewStyle().Foreground(style).Bold(true)

	valueStyle := lipgloss.Color("#b3cedc")
	desc := lipgloss.NewStyle().Foreground(valueStyle).Bold(false)
	for _, item := range pl.list {

		str := title.Render(item.title+" ") + desc.Render(item.desc)
		section = append(section, str)
	}
	return strings.Join(section, "\n")
}
