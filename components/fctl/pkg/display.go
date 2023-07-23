package fctl

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Display struct {
	header     *ui.Header
	controller config.Controller
	renderer   modelutils.Model
}

func NewDisplay() *Display {
	return &Display{
		header:     ui.NewHeader(),
		controller: nil,
		renderer:   nil,
	}
}

func (d *Display) SetHeader(model *ui.Header) *Display {
	d.header = model
	return d
}

func (d *Display) SetController(c config.Controller) *Display {
	d.controller = c
	return d
}

func (d *Display) ResetModels() *Display {
	d.controller = nil
	return d
}

func (d *Display) GetKeyMapAction() *config.KeyMapHandler {
	if d.controller != nil {
		return d.controller.GetKeyMapAction()
	}
	return nil

}

func (d *Display) Init() tea.Cmd {
	var cmd tea.Cmd = nil

	var keys = d.GetKeyMapAction()
	if keys != nil {
		d.header = d.header.AddKeyBinding(keys)
	}

	if d.controller != nil {

		renderer, err := d.controller.Run()
		if err != nil {
			panic(err)
		}

		m, err := renderer.Render()
		if err != nil {
			panic(err)
		}

		d.renderer = m
		cmd = tea.Batch(cmd, m.Init())
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

	d.header = m

	if d.controller != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			msg.Height -= 15

			m, cmd := d.renderer.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			d.renderer = m
		case tea.KeyMsg:
			m, cmd := d.renderer.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.renderer = m
		}

	}

	return d, tea.Batch(cmds...)
}

func (d *Display) View() string {
	var s = []string{
		d.header.View(),
	}

	if d.controller != nil {
		s = append(s, d.renderer.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, s...)
}
