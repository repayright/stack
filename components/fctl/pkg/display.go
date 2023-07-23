package fctl

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/formancehq/fctl/pkg/ui/theme"
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
		d.header = d.header.SetKeyBinding(keys)

	}

	cmd = d.header.Init()
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

	if d.controller != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			w, h, err := modelutils.GetTerminalSize()
			if err != nil {
				return d, tea.Quit
			}
			newMsg := tea.WindowSizeMsg{
				Width:  w,
				Height: h - d.header.GetMaxPossibleHeight() - theme.DocStyle.GetHorizontalPadding() - theme.DocStyle.GetHorizontalMargins(),
			}
			//fmt.Println("window size msg", msg.Height)
			m, cmd := d.renderer.Update(newMsg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.renderer = m
		case tea.KeyMsg:
			if keyMapHandler := d.controller.GetKeyMapAction(); keyMapHandler != nil {
				action := keyMapHandler.GetAction(tea.Key(msg))
				//fmt.Println("action", tea.Key(msg))
				if action != nil {
					controller := action(d.renderer)
					if controller == nil {
						return d, nil
					}
					d.controller = controller

					renderer, err := d.controller.Run()
					if err != nil {
						//fmt.Println("error", err)
						return d, tea.Quit
					}

					model, err := renderer.Render()
					if err != nil {
						//fmt.Println("error", err)
						return d, tea.Quit
					}

					if model == nil {
						return d, nil
					}

					d.renderer = model
					w, h, err := modelutils.GetTerminalSize()
					if err != nil {
						return d, tea.Quit
					}

					newMsg := tea.WindowSizeMsg{
						Width:  w,
						Height: h - d.header.GetMaxPossibleHeight() - 4,
					}
					m, cmd := d.renderer.Update(newMsg)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
					d.renderer = m

					// Update headers actions
					d.ResetKeyMapAction()
					header, cmd := d.header.Update(newMsg)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
					d.header = header
				}
			}

			m, cmd := d.renderer.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.renderer = m
		}

	}
	m, cmd := d.header.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	d.header = m
	return d, tea.Batch(cmds...)
}

func (d *Display) ResetKeyMapAction() *Display {
	if d.controller != nil {
		if handler := d.controller.GetKeyMapAction(); handler != nil {
			d.header.SetKeyBinding(handler)
			return d
		}
		d.header.ResetBinding()
	}
	return d
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
