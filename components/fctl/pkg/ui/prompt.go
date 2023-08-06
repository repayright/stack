package ui

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

type Prompt struct {
	style     lipgloss.Style
	model     textinput.Model
	lastInput string

	keyMapAction *config.KeyMapHandler
	yPosition    int

	commands    *Commands
	suggestions []*list.HorizontalItem
}

type Commands struct {
	commands []string
	descmap  map[string]string
}

func NewPrompt(cmd *cobra.Command) *Prompt {
	return &Prompt{
		style: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
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
	return 3
}

func (p *Prompt) GetSuggestions() []*list.HorizontalItem {
	return p.suggestions
}

func (p *Prompt) GetCursorPosition() int {
	return p.model.Position()
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
	p.model = textinput.New()
	p.model.Placeholder = "Typing..."
	p.model.CharLimit = 60

	return nil
}

func (p *Prompt) Update(msg tea.Msg) (*Prompt, tea.Cmd) {
	var cmd tea.Cmd

	p.model, cmd = p.model.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.style = p.style.Width(msg.Width - p.style.GetHorizontalBorderSize())

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			p.model.Blur()
			p.model.Reset()
			return p, func() tea.Msg {
				return modelutils.ClosePromptMsg{}
			}
		}

		v := p.model.Value()

		if p.lastInput == v {
			return p, nil
		}

		p.lastInput = v
		p.suggestions = nil

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
				item := list.NewHorizontalItem(r, p.commands.descmap[r])
				rows = append(rows, item)
			}
			p.suggestions = rows
		}()

	}

	return p, cmd
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}

func NewCommands(cmd *cobra.Command) *Commands {
	var commands []string
	var description map[string]string = make(map[string]string)
	for _, c := range cmd.Commands() {
		commands = append(commands, c.Use)
		description[c.Use] = c.Short
		if c.HasSubCommands() {
			subCommands := NewCommands(c)
			// Prefix sub tree with parent command
			for i, sub := range subCommands.commands {
				subCommands.commands[i] = c.Use + " " + sub
				description[subCommands.commands[i]] = subCommands.descmap[sub]
			}

			commands = append(commands, subCommands.commands...)
		}
	}
	return &Commands{
		commands: commands,
		descmap:  description,
	}
}
