package fctl

import (
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
	"github.com/formancehq/fctl/pkg/utils"
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

	// Get the first rendered body
	body, err := renderer.Render()
	if err != nil {
		return tea.Quit
	}
	d.renderer = body

	return tea.Sequence(
		d.header.Init(),
		d.prompt.Init(),
		d.renderer.Init(),
		d.initTermSize(),
	)
}

func (d *Display) initTermSize() tea.Cmd {
	// Init sizes
	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		return func() tea.Msg {
			return tea.QuitMsg{}
		}
	}

	d.lastTermSize = tea.WindowSizeMsg{
		Width:  w,
		Height: h,
	}
	d.newBodyWindowMsg(modelutils.ResizeMsg{
		Width:  w,
		Height: h,
	})

	return func() tea.Msg {
		d.prompt.SetWidth(d.lastBodySize.Width)
		return d.lastTermSize
	}
}

func (d *Display) HeadersHeight() int {
	if d.prompt.IsFocused() {
		return d.header.GetMaxPossibleHeight() + d.prompt.GetHeight()
	}
	return d.header.GetMaxPossibleHeight()
}

func (d *Display) newBodyWindowMsg(msg modelutils.ResizeMsg) {
	d.lastBodySize = tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: utils.Max(msg.Height-d.HeadersHeight()-4, 0), // This is due to style rendering, i need to retrive and the y space
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
		d.newBodyWindowMsg(msg)
		m, rCmd := d.renderer.Update(d.lastBodySize)
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
	d.rendered = ""
	d.addRenderHeaderView()     //1 <=> Ordered List
	d.addPromptView()           //2 <=> Maybe i can add an Renderer(map[rule]tea.Model, outString string)
	d.addRenderControllerView() //3 <=> Where rule coul
	d.addRenderSuggestionView() //4 <=> And out the output,
	d.addRenderConfirmView()    //5 <=> In some case the output is a string, a rendered view
}

func (d *Display) addPromptView() {
	if d.prompt.IsFocused() {
		d.rendered = lipgloss.JoinHorizontal(lipgloss.Top, d.rendered, d.prompt.View())
	}
}

func (d *Display) addRenderHeaderView() {
	d.rendered = lipgloss.JoinHorizontal(lipgloss.Top, d.rendered, d.header.View())
}

func (d *Display) addRenderControllerView() {
	if d.controller == nil {
		return
	}

	d.rendered = lipgloss.JoinVertical(lipgloss.Top, d.rendered, d.renderer.View())
}

func (d *Display) addRenderSuggestionView() {
	if d.suggestions == nil {
		return
	}

	d.rendered = helpers.PlaceOverlay(
		d.prompt.GetCursorPosition(),
		d.header.GetMaxPossibleHeight()+3,
		d.suggestions.View(),
		d.rendered,
		false,
	)

}

func (d *Display) addRenderConfirmView() {
	if d.confirm == nil {
		return
	}
	h := d.confirm.Styles().GetHeight()
	w := d.confirm.Styles().GetWidth()

	box := lipgloss.Place(h, w, lipgloss.Center, lipgloss.Center, d.confirm.View())

	// From Top Left Corner Always
	posX := d.lastTermSize.Width/2 - w/2
	posY := d.lastTermSize.Height/2 - h/2

	d.rendered = helpers.PlaceOverlay(posX, posY, box, d.rendered, false)
}

func (d *Display) View() string {
	return d.rendered
}
