package prompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui/list"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/hbollon/go-edlib"
	"github.com/spf13/cobra"
)

type UpdateSuggestionMsg struct{}
type Prompt struct {
	style     lipgloss.Style
	model     textinput.Model
	lastInput string

	keyMapAction *config.KeyMapHandler
	yPosition    int

	commands    *Commands
	suggestions []*list.HorizontalItem
}

func NewPrompt(cmd *cobra.Command) *Prompt {
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
		),
		yPosition: 0,
		commands:  NewCommands(cmd),
	}
}

// Border(lipgloss.NormalBorder()),
// This is the border style for the prompt
// This add a border around the prompt
// The border add 2 mores lines to the height
func (p *Prompt) GetHeight() int {
	return p.style.GetHorizontalFrameSize()
}

func (p *Prompt) GetSuggestions() []*list.HorizontalItem {
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

	case tea.KeyMsg:
		// Log := helpers.NewLogger("PROMPT")
		switch msg.String() {
		case "esc":
			p.model.Blur()
			p.model.Reset()
			p.suggestions = nil
			return p, func() tea.Msg {
				return modelutils.ClosePromptMsg{}
			}
		}

		p.model, cmd = p.model.Update(msg)
		v := p.model.Value()

		if p.lastInput == v {
			return p, nil
		}

		p.lastInput = v
		p.suggestions = nil

		// We might want to use a go routine here
		// But we need to be careful about the order of the results
		// We might want to use a channel to send the results
		ready := make(chan bool)
		go func() {
			// Cosine distance OK
			// Levenshtein distance OK
			res, err := edlib.FuzzySearchSetThreshold(v, p.commands.commands, 4, 0.3, edlib.Levenshtein)
			if err != nil || len(res) == 0 {
				return
			}

			var rows []*list.HorizontalItem = make([]*list.HorizontalItem, 0)
			for _, r := range res {
				if r == "" {
					continue
				}
				rows = append(rows, list.NewHorizontalItem(r, p.commands.descmap[r]))

			}
			ready <- true
			p.suggestions = rows
		}()
		if <-ready {
			return p, func() tea.Msg {
				return UpdateSuggestionMsg{}
			}
		} else {
			return p, func() tea.Msg {
				return UpdateSuggestionMsg{}
			}
		}

	default:
		p.model, cmd = p.model.Update(msg)
	}

	return p, cmd
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}
