package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Prompt struct {
	style lipgloss.Style
	model textinput.Model
}

func NewPrompt() *Prompt {
	return &Prompt{
		style: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
	}
}

// Border(lipgloss.NormalBorder()),
// This is the border style for the prompt
// This add a border around the prompt
// The border add 2 mores lines to the height
func (p *Prompt) GetHeight() int {
	return 3
}

func (p *Prompt) Init() tea.Cmd {
	p.model = textinput.NewModel()
	p.model.Placeholder = "Enter name"
	p.model.Focus()
	return nil
}

func (p *Prompt) Update(msg tea.Msg) (*Prompt, tea.Cmd) {
	var cmd tea.Cmd
	p.model, cmd = p.model.Update(msg)
	return p, cmd
}

func (p *Prompt) View() string {
	return p.style.Render(p.model.View())
}
