package fctl

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/header"
	"github.com/formancehq/fctl/pkg/ui/helpers"
	"github.com/formancehq/fctl/pkg/ui/list"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/formancehq/fctl/pkg/ui/prompt"
	"github.com/spf13/cobra"
)

var lock = &sync.Mutex{}
var instance *Display

type Display struct {
	header      *header.Header
	prompt      *prompt.Prompt
	suggestions *prompt.Suggestions
	controller  config.Controller
	renderer    modelutils.Model

	confirm *ui.Confirm

	lastBodySize *tea.WindowSizeMsg
	lastTermSize tea.WindowSizeMsg

	rendered string
}

func NewDisplay(cmd *cobra.Command) *Display {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &Display{
				header:     header.NewHeader(),
				controller: nil,
				renderer:   nil,
				prompt:     prompt.NewPrompt(cmd),
			}
		}
	}
	return instance
}

func (d *Display) SetHeader(model *header.Header) *Display {
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

	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		panic(err)
	}

	d.lastTermSize = tea.WindowSizeMsg{
		Width:  w,
		Height: h,
	}

	d.addControllerPromptKeyBinding(d.controller)
	d.addPromptExitKeyBinding()
	d.GenerateKeyMapAction()

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

	return tea.Sequence(
		d.header.Init(),
		d.prompt.Init(),
		d.renderer.Init(),
		func() tea.Msg {
			return d.lastTermSize
		},
	)
}

func (d *Display) HeadersHeight() int {
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
		Height: h - d.HeadersHeight(),
	}, nil
}

func (d *Display) addControllerPromptKeyBinding(c config.Controller) {
	c.GetKeyMapAction().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Turn on the prompt"),
		),
		func(model tea.Model) tea.Msg {
			if d.prompt.IsFocused() {
				return nil
			}
			return modelutils.OpenPromptMsg{}
		},
	).AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "Exit program"),
		),
		func(m tea.Model) tea.Msg {
			return tea.QuitMsg{}
		},
	)
}

func (d *Display) addPromptExitKeyBinding() {
	d.prompt.GetKeyMapAction().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Quit the prompt"),
		),
		func(model tea.Model) tea.Msg {
			return modelutils.ClosePromptMsg{}
		},
	).AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Exit program"),
		),
		func(m tea.Model) tea.Msg {
			return tea.QuitMsg{}
		},
	)
}

func (d *Display) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		header, cmd := d.header.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		d.header = header

		bodyMsg, err := d.newBodyWindowMsg()
		if err != nil {
			return d, tea.Quit
		}

		d.lastBodySize = bodyMsg
		d.lastTermSize = msg
		log := helpers.NewLogger("Display")
		log.Log("Body width", strconv.Itoa(d.lastBodySize.Width), "heigth", strconv.Itoa(d.lastBodySize.Height))
		log.Log("Term width", strconv.Itoa(d.lastTermSize.Width), "heigth", strconv.Itoa(d.lastTermSize.Height))
		m, cmd := d.renderer.Update(*bodyMsg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		d.renderer = m
		d.Render()
	case modelutils.ChangeViewMsg:
		d.controller = msg.Controller
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
		d.renderer.Init()
		d.addControllerPromptKeyBinding(d.controller)
		d.GenerateKeyMapAction()
		cmds = append(cmds, func() tea.Msg {
			return modelutils.RenderMsg{}
		})
	case modelutils.OpenPromptMsg:
		d.prompt.SwitchFocus()
		d.GenerateKeyMapAction()
		d.Render()
		cmds = append(cmds, func() tea.Msg {
			return modelutils.RenderMsg{}
		})
	case modelutils.ClosePromptMsg:
		d.GenerateKeyMapAction()
		d.suggestions = nil
		cmds = append(cmds, func() tea.Msg {
			return modelutils.RenderMsg{}
		})
	case modelutils.RenderMsg:
		d.Render()
	case prompt.UpdateSuggestionMsg:
		// Prompt Suggestions
		if !d.prompt.IsFocused() {
			d.suggestions = nil
		}
		if len(d.prompt.GetSuggestions()) > 0 {
			model := list.NewPointList(
				d.prompt.GetSuggestions()...,
			)
			d.suggestions = prompt.NewSuggestions(model)
		}
	case tea.KeyMsg:
		if d.prompt.IsFocused() {
			m, cmd := d.prompt.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.prompt = m
			cmds = append(cmds, func() tea.Msg {
				return modelutils.RenderMsg{}
			})
			break
		}

		// Check for action key binding
		// It migth change to a specific tea.Msg
		if keyMapHandler := d.controller.GetKeyMapAction(); keyMapHandler != nil {
			action := keyMapHandler.GetAction(tea.Key(msg))
			if action != nil {
				newMsg := action(d.renderer)
				if newMsg != nil {
					cmds = append(cmds, func() tea.Msg {
						return newMsg
					})
				}
			}
		}

		// Always rerender the view
		// This is needed to update the key bindings
		// All sub views are rerendered
		// FIXME: This is not efficient
		// Each vue should trigger the msg when needed
		renderer, cmd := d.renderer.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		d.renderer = renderer
		cmds = append(cmds, func() tea.Msg {
			return modelutils.RenderMsg{}
		})
	default:
		if d.prompt.IsFocused() {
			m, cmd := d.prompt.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.prompt = m
			cmds = append(cmds, func() tea.Msg {
				return modelutils.RenderMsg{}
			})
		} else {
			renderer, cmd := d.renderer.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			d.renderer = renderer
		}
	}

	return d, tea.Batch(cmds...)
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

func (d *Display) Render() {
	var s = []string{
		d.header.View(),
	}

	if d.prompt.IsFocused() {
		s = append(s, d.prompt.View())
	}

	if d.controller != nil {
		log := helpers.NewLogger("Display")
		log.Log(strconv.Itoa(d.lastTermSize.Height))
		s = append(s, d.renderer.View())
	}

	screen := lipgloss.JoinVertical(lipgloss.Top, s...)
	d.rendered = screen

	// Prompt Suggestions
	if d.suggestions != nil {
		d.rendered = helpers.PlaceOverlay(d.prompt.GetCursorPosition()+4, d.header.GetMaxPossibleHeight()+3, d.suggestions.View(), d.rendered, false)
	}

	// Confirm View
	if d.confirm == nil {
		return
	}

	confirmView := d.confirm.View()
	box := lipgloss.Place(20, 10, lipgloss.Center, lipgloss.Center, confirmView)

	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		panic(err)
	}
	str := helpers.PlaceOverlay(w/2-10, h/2-5, box, screen, false)

	d.rendered = fmt.Sprint(str)

}

func (d *Display) View() string {
	return d.rendered
}
