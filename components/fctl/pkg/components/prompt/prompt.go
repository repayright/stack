package prompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/components/table"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/modelutils"
	"github.com/hbollon/go-edlib"
)

type KeyMap struct {
	Up key.Binding
}

type UpdateSuggestionMsg struct{}
type Prompt struct {
	style lipgloss.Style
	model textinput.Model

	keyMapAction *config.KeyMapHandler
	yPosition    int

	commands *Commands

	suggestions *Suggestions
}

func NewPrompt(node *config.Node) *Prompt {
	return &Prompt{
		style: lipgloss.NewStyle().Border(lipgloss.NormalBorder()).UnsetPadding().UnsetMargins(),
		keyMapAction: config.NewKeyMapHandler().AddNewKeyBinding(
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "Validate the prompt"),
			),
			func(model tea.Model) tea.Msg {
				return nil
			},
		).AddNewKeyBinding(
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "Quit the prompt"),
			),
			func(model tea.Model) tea.Msg {
				return modelutils.ClosePromptMsg{}
			},
		),
		yPosition: 0,
		commands:  NewCommands(node),
	}
}

// Border(lipgloss.NormalBorder()),
// This is the border style for the prompt
// This add a border around the prompt
// The border add 2 mores lines to the height
func (p *Prompt) GetHeight() int {
	return p.style.GetHorizontalFrameSize()
}

func (p *Prompt) GetSuggestions() *Suggestions {
	return p.suggestions
}

func (p *Prompt) GetCursorPosition() int {
	return p.model.Position() + p.model.Cursor.Style.GetHorizontalFrameSize() + 2
}

func (p *Prompt) SetWidth(width int) *Prompt {
	p.style = p.style.Width(width - p.style.GetHorizontalFrameSize())
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
	p.model = textinput.New()
	p.model.Placeholder = "Typing..."
	p.model.CharLimit = 60
	return nil
}

func (p *Prompt) Update(msg tea.Msg) (*Prompt, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= p.style.GetHorizontalBorderSize() - p.style.GetHorizontalMargins()
		p.style = p.style.Width(
			msg.Width,
		)
		p.model, cmd = p.model.Update(msg)
		return p, cmd
	case tea.KeyMsg:
		// Log := helpers.NewLogger("PROMPT")
		p.model, cmd = p.model.Update(msg)
		v := p.model.Value()

		switch msg.String() {
		case "esc":
			p.model.Blur()
			p.model.Reset()
			p.suggestions = nil
			return p, func() tea.Msg {
				return modelutils.ClosePromptMsg{}
			}
		case "ctrl+c":
			return p, tea.Quit
		case "up", "down", "tab":
			if p.suggestions == nil {
				return p, nil
			}

			s, _ := p.suggestions.Update(msg)
			p.suggestions = &s
			return p, func() tea.Msg {
				return UpdateSuggestionMsg{}
			}
		case "enter":
			if p.suggestions == nil {
				return p, nil
			}

			s, _ := p.suggestions.Update(msg)
			p.suggestions = &s

			if row := p.suggestions.GetSelected(); row != nil {
				c0 := row.Items()[0].String()

				// We need to find the command
				// If c0 does not exist in the map it will lead to an error
				// We need to check if the command exist in the map
				controller := p.commands.controllerMap[c0]

				if controller == nil {
					return p, nil
				}

				// We need to check if the controller is a leaf
				p.model.Blur()
				p.model.Reset()
				p.suggestions = nil
				return p, func() tea.Msg {
					return modelutils.ChangeViewMsg{
						Controller: controller,
					}
				}
			}

		}

		// We might want to use a go routine here
		// But we need to be careful about the order of the results
		// We might want to use a channel to send the results
		// The model should be created and updated HERE
		ready := make(chan bool)
		go func() {
			// Cosine distance OK
			// Levenshtein distance OK
			res, _ := edlib.FuzzySearchSetThreshold(v, p.commands.commands, 4, 0.3, edlib.Levenshtein)
			var rows []*table.Row
			for _, r := range res {
				if r == "" {
					continue
				}

				cells := table.NewCells(
					table.NewCell(r,
						table.WithWidth(20),
					),
					table.NewCell(p.commands.descmap[r],
						table.WithWidth(40),
					),
				)

				row := table.NewRow(
					cells,
				)

				rows = append(rows, row)
			}
			if len(rows) == 0 {
				ready <- true
				return
			}

			p.suggestions = NewSuggestions(rows...)
			ready <- true
		}()
		if <-ready {
			return p, tea.Batch(cmd, func() tea.Msg {
				return UpdateSuggestionMsg{}
			})
		}
	default:
		p.model, cmd = p.model.Update(msg)
		return p, cmd
	}
	return p, nil
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}
