package fctl

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/helpers"
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

	confirm *ui.Confirm

	lastBodySize *tea.WindowSizeMsg
	lastTermSize *tea.WindowSizeMsg

	rendered string
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

	d.controller.GetKeyMapAction()
	d.addControllerPromptKeyBinding(d.controller)

	renderer, err := d.controller.Run()
	if err != nil {
		return tea.Quit
	}

	body, err := renderer.Render()
	if err != nil {
		return tea.Quit
	}

	d.renderer = body

	msg, err := d.newBodyWindowMsg()
	if err != nil {
		return tea.Quit
	}

	// This is needed to set the width of the prompt
	d.prompt.SetWidth(msg.Width)
	d.addPromptExitKeyBinding()

	// Generate
	d.GenerateKeyMapAction()
	return tea.Batch(
		d.header.Init(),
		d.prompt.Init(),
		body.Init(),
	)
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

func (d *Display) addControllerPromptKeyBinding(c config.Controller) {

	c.GetKeyMapAction().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Turn on the prompt"),
		),
		func(model tea.Model) config.Controller {
			if d.prompt.IsFocused() {
				return nil
			}
			d.prompt.SwitchFocus()
			d.GenerateKeyMapAction()
			d.Update(*d.lastBodySize)
			return nil
		},
	)
}

func (d *Display) addPromptExitKeyBinding() {
	d.prompt.GetKeyMapAction().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("exit"),
			key.WithHelp("exit", "quit the prompt"),
		),
		func(model tea.Model) config.Controller {
			d.prompt.SwitchFocus()
			d.GenerateKeyMapAction()
			return nil
		},
	)
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

		d.lastBodySize = bodyMsg
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

					// At this point we are sure that the model is generated
					d.renderer = model
					cmds = append(cmds, func() tea.Msg {
						d.renderer.Init()
						d.addControllerPromptKeyBinding(d.controller)
						d.GenerateKeyMapAction()
						d.Update(*d.lastBodySize)
						return nil

					})

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
			d.GenerateKeyMapAction()
			d.Update(*d.lastBodySize)
			return nil

		})
	}

	d.Render()
	return d, tea.Sequence(cmds...)
}

func (d *Display) GenerateKeyMapAction() *Display {
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

	d.header.ResetBinding()
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

func (d *Display) Render() string {
	var s = []string{
		d.header.View(),
	}

	if d.prompt.IsFocused() {
		s = append(s, d.prompt.View())
	}

	if d.controller != nil {
		s = append(s, d.renderer.View())
	}

	screen := lipgloss.JoinVertical(lipgloss.Top, s...)
	d.rendered = screen

	if d.confirm == nil {
		return d.rendered
	}

	cfimView := d.confirm.View()
	box := lipgloss.Place(20, 10, lipgloss.Center, lipgloss.Center, cfimView)

	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		panic(err)
	}
	str := helpers.PlaceOverlay(w/2-10, h/2-5, box, screen, false)

	d.rendered = fmt.Sprint(str)

	return d.rendered
}

func (d *Display) View() string {
	return d.rendered
}
