package ui

import (
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

// func GetDefaultHeaderArray() []string {

// 	var maxLength int = 0
// 	tab := strings.Split(GetDefaultFCTLASCII(fctl.Version), "\n")

// 	//Get the max length
// 	for i := 0; i < len(tab); i++ {
// 		if len(tab[i]) > maxLength {
// 			maxLength = len(tab[i])
// 		}
// 	}

// 	// Add spaces to the right
// 	for i := 0; i < len(tab); i++ {
// 		if len(tab[i]) != maxLength {
// 			tab[i] += strings.Repeat(" ", maxLength-len(tab[i]))
// 		}
// 	}
// 	return tab
// }

// For each in all children models
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

// Add a breakline to the header
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

		var out []string = []string{}
		if len(buffer) >= maxHeigth {

			out = append(out, buffer[:maxHeigth-1]...)
			// fmt.Println("out", maxHeigth)
			// fmt.Println("out", buffer[:maxHeigth])
			buffer = append(buffer, buffer[maxHeigth-1:]...)
			// fmt.Println("out", out)
		} else {
			out = buffer
			// fmt.Println("out", out)
			buffer = []string{}
		}

		if len(buffer) == 0 && len(out) == 0 {
			break
		}

		pointList := list.NewPointList(out...)

		h.AddModel(pointList)

	}
	return h
}
