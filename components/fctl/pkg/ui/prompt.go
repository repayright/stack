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
	dataset      *cobra.Command
	yPosition    int

	suggestions []*list.HorizontalItem
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
		dataset:   cmd,
		yPosition: 0,
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
	p.model = textinput.NewModel()
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
		v := p.model.Value()

		switch v {
		case "exit":
			p.model.Reset()
			p.model.Blur()
			return p, func() tea.Msg {
				return modelutils.BlurMsg{}
			}
		// case "profiles":
		// 	p.model.Reset()
		// 	p.model.Blur()
		// 	return p, func() tea.Msg {
		// 		// return modelutils.ChangeViewMsg{
		// 		// 	controller: controller,
		// 		// }
		// 		return nil
		// 	}
		default:
			if p.lastInput == v {
				return p, nil
			}

			p.lastInput = v
			p.suggestions = nil

			go func() {
				//Get map keys
				rootCmd := p.dataset.Root()
				flatCommandtree, descsMap := getTreeCommand(rootCmd)
				// Cosine distance OK
				// Levenshtein distance OK
				res, err := edlib.FuzzySearchSetThreshold(v, flatCommandtree, 4, 0.3, edlib.Cosine)
				if err != nil || len(res) == 0 {
					return
				}

				var rows []*list.HorizontalItem = make([]*list.HorizontalItem, 0)
				for _, r := range res {
					if r == "" {
						continue
					}
					item := list.NewHorizontalItem(r, descsMap[r])
					rows = append(rows, item)
				}
				p.suggestions = rows
			}()

		}
	}

	return p, cmd
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}

func getTreeCommand(cmd *cobra.Command) ([]string, map[string]string) {
	var commands []string
	var description map[string]string = make(map[string]string)
	for _, c := range cmd.Commands() {
		commands = append(commands, c.Use)
		description[c.Use] = c.Short
		if c.HasSubCommands() {
			subTree, desc := getTreeCommand(c)
			// Prefix sub tree with parent command
			for i, sub := range subTree {
				subTree[i] = c.Use + " " + sub
				description[subTree[i]] = desc[sub]
			}

			commands = append(commands, subTree...)
		}
	}
	return commands, description
}
