package fctl

import (
	"sync"

	"github.com/charmbracelet/bubbles/key"
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

	lastTerminalSize *tea.WindowSizeMsg
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

func (d *Display) Init() tea.Cmd {

	flags := d.controller.GetConfig().GetAllFLags()
	conf, err := GetConfig(flags)
	if err != nil {
		panic(err)
	}

	profiles := GetCurrentProfile(flags, conf)
	d.header.GetContext().SetFctlVersion(Version)
	d.header.GetContext().SetOrg(profiles.defaultOrganization)
	d.header.GetContext().SetProfile(conf.currentProfile)

	var keys = d.GetKeyMapAction()
	if keys != nil {
		d.header = d.header.SetKeyBinding(keys)

	}

	cmd := d.header.Init()
	d.controller.GetKeyMapAction().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Turn on the prompt"),
		),
		func(model tea.Model) config.Controller {
			if d.prompt.IsFocused() {
				return nil
			}
			d.prompt.SwitchFocus()
			d.ResetKeyMapAction()
			d.Update(*d.lastTerminalSize)
			return nil
		},
	)
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

	pCmd := d.prompt.Init()
	if pCmd != nil {
		cmd = tea.Batch(cmd, pCmd)
	}

	msg, err := d.newBodyWindowMsg()
	if err != nil {
		panic(err)
	}

	d.prompt.SetWidth(msg.Width)
	d.prompt.GetKeyMapAction().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("exit"),
			key.WithHelp("exit", "Quit the prompt"),
		),
		func(model tea.Model) config.Controller {
			d.prompt.SwitchFocus()
			d.ResetKeyMapAction()
			return nil
		},
	)

	d.ResetKeyMapAction()
	return cmd
}

func (d *Display) GetChildrenHeight() int {
	if d.prompt.IsFocused() {
		return d.header.GetMaxPossibleHeight() + d.prompt.GetHeight()

	}
	return d.header.GetMaxPossibleHeight()

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
	p, cmd := d.prompt.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	d.prompt = p
	m, cmd := d.header.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	d.header = m
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		bodyMsg, err := d.newBodyWindowMsg()
		if err != nil {
			return d, tea.Quit
		}

		d.lastTerminalSize = bodyMsg
		m, cmd := d.renderer.Update(*bodyMsg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		d.renderer = m

	case tea.KeyMsg:
		if keyMapHandler := d.controller.GetKeyMapAction(); !d.prompt.IsFocused() && keyMapHandler != nil {
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
					newMsg, err := d.newBodyWindowMsg()
					if err != nil {
						return d, nil
					}

					m, cmd := d.renderer.Update(newMsg)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
					d.renderer = m
					d.ResetKeyMapAction()

				}
			}
		}
		m, cmd := d.renderer.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		d.renderer = m
	case modelutils.BlurMsg:
		cmds = append(cmds, func() tea.Msg {
			d.ResetKeyMapAction()
			d.Update(*d.lastTerminalSize)
			msg, err := d.newBodyWindowMsg()
			if err != nil {
				return tea.Quit
			}

			return msg

		})
	}

	return d, tea.Sequence(cmds...)
}

func (d *Display) ResetKeyMapAction() *Display {
	d.header.ResetBinding()
	if d.prompt.IsFocused() {
		if handler := d.prompt.GetKeyMapAction(); handler != nil {
			d.header.SetKeyBinding(handler)
			return d
		}
	}

	if d.controller != nil {
		if handler := d.controller.GetKeyMapAction(); handler != nil {
			d.header.SetKeyBinding(handler)
			return d
		}

	}

	return d
}

func (d *Display) GetKeyMapAction() *config.KeyMapHandler {
	if d.prompt.IsFocused() {
		return d.prompt.GetKeyMapAction()
	}
	if d.controller != nil {
		return d.controller.GetKeyMapAction()
	}

	return nil

}

func (d *Display) View() string {
	var s = []string{
		d.header.View(),
	}

	if d.prompt.IsFocused() {
		s = append(s, d.prompt.View())
	}

	if d.controller != nil {
		s = append(s, d.renderer.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, s...)
}
