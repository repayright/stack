package ui

import (
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/list"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Header struct {
	version *FctlModel
	actions []list.PointList
}

func NewHeader() *Header {
	return &Header{
		actions: []list.PointList{},
		version: NewFctlModel(),
	}
}

func (h *Header) AddModel(model *list.PointList) *Header {
	h.actions = append(h.actions, *model)
	return h
}

// Depends on the FCTL version components
func (h *Header) GetMaxPossibleHeight() int {
	return h.version.GetMaxPossibleHeight()
}

func (h *Header) Init() tea.Cmd {
	return h.version.Init()
}
func (h *Header) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}
func (h *Header) Update(msg tea.Msg) (*Header, tea.Cmd) {
	return h, nil
}

func (h *Header) View() string {

	var m []string = []string{
		h.version.View(),
	}
	for _, model := range h.actions {
		m = append(m, model.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, m...)
}

func (h *Header) AddKeyBinding(keys ...*modelutils.KeyMapHandler) *Header {
	var maxHeigth int = h.GetMaxPossibleHeight()
	var buffer []string = []string{}

	for _, key := range keys {
		v := key.View()
		s := strings.Split(v, "\n")
		buffer = append(buffer, s...)
	}

	// Get the number of list to create
	n := math.Ceil(float64(len((buffer))) / float64(maxHeigth))
	i := int(n)

	var out [][]string = make([][]string, i)

	// Prepare each list with strings items
	for c := 0; c < i; c++ {

		min := c * maxHeigth
		max := (c + 1) * maxHeigth

		if max >= len(buffer) {
			max = len(buffer) - 1
		}

		out = append(out, buffer[min:max])

	}

	// Create and add models
	for _, o := range out {
		pointList := list.NewPointList(o...)
		h.AddModel(pointList)
	}

	return h
}
