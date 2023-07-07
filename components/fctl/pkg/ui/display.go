package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Display struct {
	header Header
	body   *Model

	time time.Time
}

func NewDisplay() *Display {
	return &Display{
		body: nil,
	}
}

func (d *Display) SetHeader(model *Header) *Display {
	d.header = *model
	return d
}

func (d *Display) SetCurrentModel(model *Model) *Display {
	d.body = model
	return d
}

func (d *Display) ResetModels() *Display {
	d.body = nil
	return d
}

func (d *Display) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	var keys *modelutils.KeyMapHandler
	if d.body != nil {
		keys = (*d.body).GetListKeyMapHandler()
		return keys
	}
	return nil

}

func (d *Display) Init() tea.Cmd {
	var cmd tea.Cmd = nil
	d.time = time.Now()

	// Init modelKeys
	var keys *modelutils.KeyMapHandler = d.GetListKeyMapHandler()
	d.header = *d.header.AddKeyBinding(keys)

	if d.body != nil {
		cmd = tea.Batch(cmd, (*d.body).Init())
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

	if d.body != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			msg.Height -= 15
			m, cmd := (*d.body).Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			*d.body = m
		case tea.KeyMsg:
			m, cmd := (*d.body).Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			*d.body = m
		}

	}

	return d, tea.Batch(cmds...)
}

func (d *Display) View() string {
	var s []string = []string{
		d.header.View(),
	}

	if d.body != nil {
		s = append(s, (*d.body).View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, s...)
}
