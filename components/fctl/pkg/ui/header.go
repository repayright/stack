package ui

import (
	"math"
	"strings"

	blist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/list"

	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Header struct {
	logo        *Logo
	fctlContext *Context
	modelAction []list.PointList
	rendered    string
}

func NewHeader() *Header {
	return &Header{
		modelAction: []list.PointList{},
		logo:        NewLogo(),
		fctlContext: NewContext(),
	}
}

func (h *Header) AddModel(model *list.PointList) *Header {
	h.modelAction = append(h.modelAction, *model)
	return h
}

// Depends on the FCTL version components
func (h *Header) GetMaxPossibleHeight() int {
	return h.logo.GetMaxPossibleHeight()
}

func (h *Header) Init() tea.Cmd {
	return h.logo.Init()
}
func (h *Header) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}
func (h *Header) Update(msg tea.Msg) (*Header, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		logoView := h.logo.View()
		var m []string = []string{}

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
			if i == len(h.modelAction)-1 {
				for k, line := range split {
					split[k] = lipgloss.NewStyle().PaddingLeft(1).Render(line)
				}
			}

			out := lipgloss.JoinVertical(lipgloss.Left, split...)
			m = append(m, out)
		}

		t := lipgloss.NewStyle().Render(lipgloss.JoinHorizontal(lipgloss.Top, m...))

		terminalW := msg.Width - h.logo.GetMaxPossibleWidth() - h.fctlContext.GetMaxPossibleWidth()

		leftDiv := h.fctlContext.View()

		style := lipgloss.NewStyle().PaddingTop(1)

		rightDiv := lipgloss.Place(h.logo.GetMaxPossibleWidth(), h.logo.GetMaxPossibleHeight(), lipgloss.Top, lipgloss.Top, style.Render(logoView))
		middleDiv := lipgloss.Place(terminalW, h.GetMaxPossibleHeight(), lipgloss.Center, lipgloss.Top, t)

		h.rendered = lipgloss.JoinHorizontal(lipgloss.Top, leftDiv, middleDiv, rightDiv)
	}

	return h, nil
}

func (h *Header) View() string {
	return h.rendered
}

func (h *Header) AddKeyBinding(keys ...*modelutils.KeyMapHandler) *Header {
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
	var out [][]blist.Item
	for c := 0; c < i; c++ {

		min := c * maxHeigth
		max := (c + 1) * maxHeigth

		if max >= len(buffer) {
			max = len(buffer) - 1
		}

		bloc := buffer[min:max]
		maxWidth = modelutils.GetMaxCharPosXinCharList(bloc, ":") + 1
		items := []blist.Item{}

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
	for _, o := range out {
		pl := list.NewPointList(o,
			list.NewHorizontalItemDelegate(),
			100,
			maxHeigth+1,
		)

		h.AddModel(pl)
	}

	return h
}
