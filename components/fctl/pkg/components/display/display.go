package display

import (
	"sync"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/components/header"
	"github.com/formancehq/fctl/pkg/components/prompt"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/helpers"
	"github.com/formancehq/fctl/pkg/modelutils"
	"github.com/formancehq/fctl/pkg/utils"
)

var lock = &sync.Mutex{}
var instance *Display

type Display struct {
	header      *header.Header
	prompt      *prompt.Prompt
	suggestions *prompt.Suggestions
	controller  config.Controller
	renderer    modelutils.Model

	confirm *Confirm

	lastBodySize tea.WindowSizeMsg
	lastTermSize tea.WindowSizeMsg

	rendered string

	//Tweak
	knownControllers []config.Controller
}

type DisplayOpts func(*Display) *Display

func WithController(controller config.Controller) DisplayOpts {
	return func(d *Display) *Display {
		d.controller = controller
		d.knownControllers = append(d.knownControllers, controller)
		return d
	}
}

func WithPrompt(prompt *prompt.Prompt) DisplayOpts {
	return func(d *Display) *Display {
		d.prompt = prompt
		return d
	}
}

func NewDisplay(opts ...DisplayOpts) *Display {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &Display{
				header: header.NewHeader(),
			}
		}
	}

	for _, opt := range opts {
		instance = opt(instance)
	}

	return instance
}

func (d *Display) SetHeader(model *header.Header) *Display {
	d.header = model
	return d
}

func (d *Display) ResetModels() *Display {
	d.controller = nil
	return d
}

func (d *Display) exitBinding() *config.KeyMapHandler {
	return config.NewKeyMapHandler().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Exit program"),
		),
		func(m tea.Model) tea.Msg {
			return tea.QuitMsg{}
		},
	)
}

func (d *Display) getPromptOnOpen() *config.KeyMapHandler {
	return config.NewKeyMapHandler().AddNewKeyBinding(
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
	)
}

func (d *Display) Init() tea.Cmd {
	flags := d.controller.GetConfig().GetAllFLags()
	conf, err := config.GetConfig(flags)
	if err != nil {
		panic(err)
	}
	profiles := config.GetCurrentProfile(flags, conf)
	d.header.GetContext().SetFctlVersion(config.Version)
	d.header.GetContext().SetOrg(profiles.DefaultOrganization())
	d.header.GetContext().SetProfile(conf.GetCurrentProfileName())

	d.GenerateKeyMapActionHelper()

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
	)
}

func (d *Display) HeadersHeight() int {
	if d.prompt.IsFocused() {
		return d.header.GetMaxPossibleHeight() + d.prompt.GetHeight()
	}
	return d.header.GetMaxPossibleHeight() - 1
}

func (d *Display) newBodyWindowMsg(msg modelutils.ResizeMsg) {
	d.lastBodySize = tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: utils.Max(msg.Height-d.HeadersHeight()-4, 0), // This is due to style rendering, i need to retrive and the y space from the body styles
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
		d.newBodyWindowMsg(modelutils.ResizeMsg{
			Width:  msg.Width,
			Height: msg.Height,
		})
		return d, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  msg.Width,
				Height: msg.Height,
			}
		}
	case modelutils.ChangeViewMsg:
		Log := helpers.NewLogger("ENTER")

		d.controller = msg.Controller

		renderer, err := d.controller.Run()
		if err != nil {
			Log.Logf("selected stack id: %s", err)
			return d, tea.Quit
		}

		model, err := renderer.Render()
		if err != nil {
			return d, tea.Quit
		}
		d.renderer = model

		d.GenerateKeyMapActionHelper()
		d.suggestions = nil

		cmd := d.renderer.Init()

		r, rCmd := d.renderer.Update(d.lastBodySize)
		d.renderer = r
		return d, tea.Sequence(cmd, rCmd, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  d.lastTermSize.Width,
				Height: d.lastTermSize.Height,
			}
		})

	case modelutils.ConfirmActionMsg:
		d.confirm = NewConfirm(msg)

	case modelutils.OpenPromptMsg:
		d.prompt.SwitchFocus()
		d.GenerateKeyMapActionHelper()
		return d, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  d.lastTermSize.Width,
				Height: d.lastTermSize.Height,
			}
		}
	case modelutils.ClosePromptMsg:
		d.GenerateKeyMapActionHelper()
		d.suggestions = nil
		return d, func() tea.Msg {
			return modelutils.ResizeMsg{
				Width:  d.lastTermSize.Width,
				Height: d.lastTermSize.Height,
			}
		}
	case prompt.UpdateSuggestionMsg:
		// If the prompt is not focused, reset the suggestions
		if !d.prompt.IsFocused() {
			d.suggestions = nil
		}

		// Reset model if none is provided
		d.suggestions = d.prompt.GetSuggestions()
	case modelutils.CloseConfirmMsg:
		d.confirm = nil
		return d, func() tea.Msg {
			return modelutils.RenderMsg{}
		}
	case tea.KeyMsg:
		if d.confirm != nil {
			confirm, cmd := d.confirm.Update(msg)
			d.confirm = &confirm
			return d, tea.Sequence(cmd, func() tea.Msg {
				return modelutils.RenderMsg{}
			})
		}
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

		keyMapHandler := d.getContextualMapAction()
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

func (d *Display) getContextualMapAction() *config.KeyMapHandler {
	var keyMapHandler *config.KeyMapHandler = config.NewKeyMapHandler()

	if d.prompt != nil {
		if d.prompt.IsFocused() {
			if handler := d.prompt.GetKeyMapAction(); handler != nil {
				h := handler.Copy()
				h.AddKeyMapHandler(d.exitBinding())
				return h
			}
		}
		keyMapHandler = keyMapHandler.AddKeyMapHandler(d.getPromptOnOpen())
	}

	if d.controller != nil {
		if handler := d.controller.GetKeyMapAction(); handler != nil {
			keyMapHandler = keyMapHandler.AddKeyMapHandler(handler)
		}
	}

	keyMapHandler.AddKeyMapHandler(d.exitBinding())
	return keyMapHandler
}

func (d *Display) GenerateKeyMapActionHelper() {
	d.header.ResetBinding()
	var keyMapHandler *config.KeyMapHandler = d.getContextualMapAction()
	d.header.SetKeyBinding(*keyMapHandler)
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

// Only rerender necessary part, need to keep old state
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
		d.rendered = lipgloss.JoinVertical(lipgloss.Top, d.rendered, d.prompt.View())
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

	if d.suggestions.IsEmpty() {
		return
	}

	d.rendered = helpers.PlaceOverlay(
		d.prompt.GetCursorPosition()+2,
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

	// From Top Left Corner Always
	posX := d.lastTermSize.Width/2 - w/2
	posY := d.lastTermSize.Height/2 - h/2

	d.rendered = helpers.PlaceOverlay(posX, posY, d.confirm.View(), d.rendered, false)
}

func (d *Display) View() string {
	return d.rendered
}
