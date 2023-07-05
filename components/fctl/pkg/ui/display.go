package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Display struct {
	header Header
	list   []Model

	time time.Time
}

func NewDisplay() *Display {
	return &Display{
		list: []Model{},
	}
}

func (d *Display) SetHeader(model *Header) *Display {
	d.header = *model
	return d
}

func (d *Display) AppendModels(models ...Model) *Display {
	d.list = append(d.list, models...)
	return d
}

func (d *Display) ResetModels() *Display {
	d.list = []Model{}
	return d
}

func (d *Display) GetListKeyMapHandler() []*modelutils.KeyMapHandler {
	var keys []*modelutils.KeyMapHandler
	for _, model := range d.list {
		if c := model.GetListKeyMapHandler(); c != nil {
			keys = append(keys, c)
		}
	}
	return keys

}

func (d *Display) Init() tea.Cmd {
	var cmd tea.Cmd = nil
	d.time = time.Now()

	d.header.Init()

	// Init modelKeys
	var keys []*modelutils.KeyMapHandler = d.GetListKeyMapHandler()
	d.header = *d.header.AddKeyBinding(keys...)

	// Init models
	for _, model := range d.list {
		if c := model.Init(); c != nil {
			cmd = tea.Batch(cmd, c)
		}
	}
	return cmd
}

func (d *Display) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	m, cmd := d.header.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	d.header = *m

	for k, model := range d.list {
		nm, cmd := model.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		d.list[k] = nm
	}

	return d, tea.Batch(cmds...)
}

func (d *Display) View() string {
	var s []string = []string{
		d.header.View(),
	}
	for _, model := range d.list {
		s = append(s, model.View())
	}
	return lipgloss.JoinVertical(lipgloss.Top, s...)
	// return ""
}
