package ui

import (
	"math"
	"strings"

	"github.com/formancehq/fctl/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/list"

	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Header struct {
	logo        *Logo
	fctlContext *Context
	modelAction []*list.PointList
	rendered    string
}

func NewHeader() *Header {
	return &Header{
		modelAction: make([]*list.PointList, 0),
		logo:        NewLogo(),
		fctlContext: NewContext(),
	}
}

func (h *Header) GetContext() *Context {
	return h.fctlContext
}

func (h *Header) AddModel(model *list.PointList) *Header {
	h.modelAction = append(h.modelAction, model)
	return h
}

func (h *Header) GetMaxPossibleHeight() int {
	return h.logo.GetMaxPossibleHeight()
}

func (h *Header) Init() tea.Cmd {
	return tea.Batch(
		h.fctlContext.Init(),
		h.logo.Init(),
	)
}

func (h *Header) Update(msg tea.Msg) (*Header, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerStyle := lipgloss.NewStyle().Margin(1, 1, 1, 1)
		middleDivSize := msg.Width - h.logo.GetMaxPossibleWidth() - h.fctlContext.GetMaxPossibleWidth() - 2
		rightDiv := lipgloss.Place(h.logo.GetMaxPossibleWidth(), h.logo.GetMaxPossibleHeight(), lipgloss.Top, lipgloss.Top, h.logo.View())
		leftDiv := lipgloss.Place(h.fctlContext.GetMaxPossibleWidth(), h.fctlContext.GetMaxPossibleHeight(), lipgloss.Top, lipgloss.Top, h.fctlContext.View())

		if len(h.modelAction) == 0 {
			h.rendered = headerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, leftDiv, strings.Repeat(" ", middleDivSize), rightDiv))
			return h, nil
		}

		var m = make([]string, 0)
		for i, model := range h.modelAction {
			tmp := model.View()
			split := strings.Split(tmp, "\n")
			// First add padding right
			if i == 0 {
				for k, line := range split {
					split[k] = lipgloss.NewStyle().PaddingRight(1).Render(line)
				}

			}
			// Last add padding left
			if i == len(h.modelAction)-1 && i > 0 {
				for k, line := range split {
					split[k] = lipgloss.NewStyle().PaddingLeft(1).Render(line)
				}
			}

			out := lipgloss.JoinVertical(lipgloss.Left, split...)
			m = append(m, out)
		}

		text := lipgloss.NewStyle().Render(lipgloss.JoinHorizontal(lipgloss.Top, m...))
		middleDiv := lipgloss.Place(middleDivSize, h.GetMaxPossibleHeight(), lipgloss.Center, lipgloss.Top, text)
		h.rendered = headerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, leftDiv, middleDiv, rightDiv))

		return h, nil
	}

	return h, nil
}

func (h *Header) View() string {
	return h.rendered
}

func (h *Header) GetListKeyMapHandler() *config.KeyMapHandler {
	return nil
}
func (h *Header) ResetBinding() *Header {
	h.modelAction = make([]*list.PointList, 0)
	return h
}

func (h *Header) SetKeyBinding(keys ...*config.KeyMapHandler) *Header {
	var maxHeigth int = h.GetMaxPossibleHeight()
	var maxWidth int
	var buffer []string = []string{}

	for _, key := range keys {
		v := key.View()
		s := strings.Split(v, "\n")
		buffer = append(buffer, s...)
	}

	// Get the number of list to create
	n := math.Ceil(float64(len((buffer))) / float64(maxHeigth))
	i := int(n)
	var out = make([][]*list.HorizontalItem, 0)
	for c := 0; c < i; c++ {

		min := c * maxHeigth
		max := (c + 1) * maxHeigth

		if max >= len(buffer) {
			max = len(buffer) - 1
		}

		bloc := buffer[min:max]
		maxWidth = modelutils.GetMaxCharPosXinCharList(bloc, ":") + 1
		var items = make([]*list.HorizontalItem, 0)
		for _, v := range bloc {

			split := strings.Split(v, ":")

			part0 := modelutils.FillCharBeforeChar(split[0]+" :", " ", ":", maxWidth)
			l := list.NewHorizontalItem(part0, strings.TrimPrefix(split[1], " "))
			items = append(items, l)
		}

		if len(items) > 0 {
			out = append(out, items)
		}

	}

	h.modelAction = make([]*list.PointList, 0)

	for _, o := range out {
		pl := list.NewPointList(o...)
		h.AddModel(pl)
	}

	return h
}
