package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Display[T any] struct {
	header       Header
	body         modelutils.Controller[T]
	bodyRenderer modelutils.Model
	time         time.Time
}

func NewDisplay[T any]() *Display[T] {
	return &Display[T]{
		body: nil,
	}
}

func (d *Display[T]) SetHeader(model *Header) *Display[T] {
	d.header = *model
	return d
}

func (d *Display[T]) SetBody(c modelutils.Controller[T]) *Display[T] {
	d.body = c
	return d
}

func (d *Display[T]) ResetModels() *Display[T] {
	d.body = nil
	return d
}

func (d *Display[T]) GetKeyMapAction() *modelutils.KeyMapHandler[T] {
	if d.body != nil {
		return d.body.GetKeyMapAction()
	}
	return nil

}

func (d *Display[T]) Init() tea.Cmd {
	var cmd tea.Cmd = nil
	d.time = time.Now()

	var keys *modelutils.KeyMapHandler[T] = d.GetKeyMapAction()
	if keys != nil {
		d.header = *d.header.AddKeyBinding(keys)
	}

	if d.body != nil {

		m, err := d.body.Render()
		if err != nil {
			panic(err)
		}

		d.bodyRenderer = m
		cmd = tea.Batch(cmd, m.Init())
	}

	return cmd
}

func (d *Display[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			m, cmd := d.bodyRenderer.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			d.bodyRenderer = m
		case tea.KeyMsg:
			m, cmd := d.bodyRenderer.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.bodyRenderer = m
		}

	}

	return d, tea.Batch(cmds...)
}

func (d *Display[T]) View() string {
	var s []string = []string{
		//d.header.View(),
	}

	if d.body != nil {
		s = append(s, d.bodyRenderer.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, s...)
}
