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
		//Tabulation
		case "tab":
			if p.suggestions == nil {
				return p, nil
			}

			// if p.cursorY+1 > len-1 {
			// 	p.cursorY = 0
			// } else {
			// 	p.cursorY++
			// }
		case "enter":
		}

		// if p.lastInput == v {
		// 	return p, nil
		// }

		// p.lastInput = v
		p.suggestions = nil

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
						table.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Background(lipgloss.Color("#000000"))),
						table.WithWidth(20),
					),
					table.NewCell(p.commands.descmap[r],
						table.WithStyle(
							lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Background(lipgloss.Color("#000000")),
						),
						table.WithWidth(40),
					),
				)

				row := table.NewRow(
					cells,
					table.WithRowStyle(
						lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Background(lipgloss.Color("#000000")),
					),
				)

				rows = append(rows, row)
			}
			if len(rows) >= 0 {
				p.suggestions = NewSuggestions(rows...)
			}
			ready <- true
		}()
		if <-ready {
			return p, func() tea.Msg {
				return UpdateSuggestionMsg{}
			}
		}
	default:
		p.model, cmd = p.model.Update(msg)
		return p, cmd
	}

	return p, cmd
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}
