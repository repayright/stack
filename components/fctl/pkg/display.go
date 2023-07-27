package fctl

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

var lock = &sync.Mutex{}
var instance *Display

type Display struct {
	header     *ui.Header
	prompt     *ui.Prompt
	controller config.Controller
	renderer   modelutils.Model
}

func NewDisplay() *Display {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &Display{
				header:     ui.NewHeader(),
				controller: nil,
				renderer:   nil,
				prompt:     ui.NewPrompt(),
			}
		}
	}
	return instance
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
	flags := d.controller.GetConfig().GetAllFLags()
	config, err := GetConfig(flags)
	if err != nil {
		panic(err)
	}

	profiles := GetCurrentProfile(flags, config)
	d.header.GetContext().SetFctlVersion(Version)
	d.header.GetContext().SetOrg(profiles.defaultOrganization)
	d.header.GetContext().SetProfile(config.currentProfile)

	var keys = d.GetKeyMapAction()
	if keys != nil {
		d.header = d.header.SetKeyBinding(keys)

	}

	cmd := d.header.Init()

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

	return tea.Batch(cmd, d.prompt.Init())
}

func (d *Display) GetChildrenHeight() int {
	return d.header.GetMaxPossibleHeight() + d.prompt.GetHeight()
}

func (d *Display) newBodyWindowMsg() (*tea.WindowSizeMsg, error) {
	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		return nil, err
	}
	return &tea.WindowSizeMsg{
		Width:  w,
		Height: h - d.GetChildrenHeight() - theme.DocStyle.GetHorizontalPadding() - theme.DocStyle.GetHorizontalMargins(),
	}, nil
}

func (d *Display) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	if d.controller != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:

			bodyMsg, err := d.newBodyWindowMsg()
			if err != nil {
				return d, tea.Quit
			}
			m, cmd := d.renderer.Update(*bodyMsg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.renderer = m

			p, cmd := d.prompt.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			d.prompt = p

		case tea.KeyMsg:
			if keyMapHandler := d.controller.GetKeyMapAction(); keyMapHandler != nil {
				action := keyMapHandler.GetAction(tea.Key(msg))
				if action != nil {
					controller := action(d.renderer)
					if controller != nil {
						d.controller = controller
						renderer, err := d.controller.Run()
						if err != nil {
							return d, tea.Quit
						}

						model, err := renderer.Render()
						if err != nil {
							return d, tea.Quit
						}

						// A controller keymap action does not always return a model
						if model == nil {
							return d, nil
						}

						// At this point we are sure that the model is generated
						d.renderer = model
						w, h, err := modelutils.GetTerminalSize()
						if err != nil {
							return d, tea.Quit
						}

						newMsg := tea.WindowSizeMsg{
							Width:  w,
							Height: h - d.GetChildrenHeight() - theme.DocStyle.GetHorizontalPadding() - theme.DocStyle.GetHorizontalMargins(),
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

	p, cmd := d.prompt.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	d.prompt = p
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
		d.prompt.View(),
	}

	if d.controller != nil {
		s = append(s, d.renderer.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, s...)
}
