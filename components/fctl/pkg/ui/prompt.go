package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Prompt struct {
	style lipgloss.Style
	model textinput.Model

	keyMapAction *config.KeyMapHandler
}

func NewPrompt() *Prompt {
	return &Prompt{
		style: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		keyMapAction: config.NewKeyMapHandler().AddNewKeyBinding(
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "Validate the prompt"),
			),
			func(model tea.Model) config.Controller {
				return nil
			},
		),
	}
}

// Border(lipgloss.NormalBorder()),
// This is the border style for the prompt
// This add a border around the prompt
// The border add 2 mores lines to the height
func (p *Prompt) GetHeight() int {
	return 3
}

func (p *Prompt) SetWidth(width int) *Prompt {
	p.style = p.style.Width(width)
	return p
}

func (p *Prompt) SwitchFocus() tea.Cmd {
	return p.model.Focus()
}

func (p *Prompt) IsFocused() bool {
	return p.model.Focused()
}

func (p *Prompt) GetKeyMapAction() *config.KeyMapHandler {
	return p.keyMapAction
}

func (p *Prompt) Init() tea.Cmd {
	p.model = textinput.NewModel()
	p.model.Placeholder = "Typing..."
	return nil
}

func (p *Prompt) Update(msg tea.Msg) (*Prompt, tea.Cmd) {
	var cmd tea.Cmd
	p.model, cmd = p.model.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.style = p.style.Width(msg.Width - p.style.GetHorizontalBorderSize())

	case tea.KeyMsg:
		v := p.model.Value()
		switch v {
		case "exit":
			p.model.Reset()
			p.model.Blur()
			return p, func() tea.Msg {
				return modelutils.BlurMsg{}
			}
		case "profiles":
			p.model.Reset()
			p.model.Blur()
			// config := profiles.NewListConfig()
			// config.SetOut(config.GetOut())
			// config.SetContext(config.GetContext())
			// controller := profiles.NewListController(config)
			return p, func() tea.Msg {
				// return modelutils.ChangeViewMsg{
				// 	controller: controller,
				// }
				return nil
			}
		}
	}

	return p, cmd
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}
