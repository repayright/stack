package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type PointList struct {
	points []string
}

func NewPointList(p ...string) *PointList {
	return &PointList{
		points: p,
	}
}
func (pl *PointList) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}
func (pl *PointList) Init() tea.Cmd {
	return nil
}
func (pl *PointList) Update(msg tea.Msg) (*PointList, tea.Cmd) {

	return pl, nil
}

func (pl *PointList) View() string {

	var str string = ""
	for _, point := range pl.points {
		str += "• " + point + "\n"
	}

	return str
}
