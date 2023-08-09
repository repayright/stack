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

	lastBodySize tea.WindowSizeMsg
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

func (d *Display) Init() tea.Cmd {

	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		panic(err)
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

	w, h, err = modelutils.GetTerminalSize()
	if err != nil {
		return tea.Quit
	}

	msg := d.newBodyWindowMsg(modelutils.ResizeMsg{
		Width:  w,
		Height: h,
	})
	if err != nil {
		return tea.Quit
	}

	// This is needed to set the width of the prompt
	d.prompt.SetWidth(msg.Width)
	d.lastTermSize = tea.WindowSizeMsg{
		Width:  w,
		Height: h,
	}
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
		Log := helpers.NewLogger("Display")
		Log.Log(fmt.Sprintf("Header Height: %d", d.header.GetMaxPossibleHeight()+d.prompt.GetHeight()+1))
		return d.header.GetMaxPossibleHeight() + d.prompt.GetHeight() + 1

	}
	return d.header.GetMaxPossibleHeight()

}

func (d *Display) newBodyWindowMsg(msg modelutils.ResizeMsg) tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: msg.Height - d.HeadersHeight(),
	}
}

func (d *Display) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case modelutils.RenderMsg:
		d.Render()
		return d, nil
	case modelutils.ResizeMsg:
		header, cmd := d.header.Update(d.lastTermSize)
		d.header = header
		d.prompt.SetWidth(msg.Width)

		bodyMsg := d.newBodyWindowMsg(msg)
		d.lastBodySize = bodyMsg
		log := helpers.NewLogger("Display")
		log.Log("Body width", strconv.Itoa(d.lastBodySize.Width), "heigth", strconv.Itoa(d.lastBodySize.Height))
		log.Log("Term width", strconv.Itoa(d.lastTermSize.Width), "heigth", strconv.Itoa(d.lastTermSize.Height))
		m, rCmd := d.renderer.Update(bodyMsg)
		d.renderer = m
		return d, tea.Sequence(cmd, rCmd, func() tea.Msg {
			return modelutils.RenderMsg{}
		})
	case tea.WindowSizeMsg:
		d.lastTermSize = msg

		return d, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  msg.Width,
				Height: msg.Height,
			}
		}
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

		d.addControllerPromptKeyBinding(d.controller)
		d.GenerateKeyMapAction()
		return d, tea.Sequence(
			d.renderer.Init(),
			func() tea.Msg {
				return modelutils.RenderMsg{}
			})
	case modelutils.OpenPromptMsg:
		d.prompt.SwitchFocus()
		d.GenerateKeyMapAction()
		return d, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  d.lastTermSize.Width,
				Height: d.lastTermSize.Height,
			}
		}
	case modelutils.ClosePromptMsg:
		d.GenerateKeyMapAction()
		d.suggestions = nil
		return d, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  d.lastTermSize.Width,
				Height: d.lastTermSize.Height,
			}
		}
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
			d.prompt = m
			return d, tea.Sequence(cmd, func() tea.Msg {
				return modelutils.RenderMsg{}
			})
		}
		renderer, cmd := d.renderer.Update(msg)
		d.renderer = renderer

		// Check for action key binding
		// It migth change to a specific tea.Msg
		var postCmd tea.Cmd = tea.Sequence(cmd, func() tea.Msg {
			return modelutils.RenderMsg{}
		})

		keyMapHandler := d.controller.GetKeyMapAction()

		// No key map handler
		if keyMapHandler == nil {
			return d, postCmd
		}

		action := keyMapHandler.GetAction(tea.Key(msg))

		// No real action defined
		if action == nil {
			return d, postCmd
		}

		newMsg := action(d.renderer)

		// No new message
		if newMsg == nil {
			return d, postCmd
		}

		// New message
		postCmd = tea.Sequence(func() tea.Msg {
			return newMsg
		}, postCmd)

		return d, postCmd
	default:
		if d.prompt.IsFocused() {
			m, cmd := d.prompt.Update(msg)
			d.prompt = m
			return d, tea.Sequence(cmd, func() tea.Msg {
				return modelutils.RenderMsg{}
			})
		} else {
			renderer, cmd := d.renderer.Update(msg)
			d.renderer = renderer
			return d, cmd
		}
	}

	return d, nil
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
