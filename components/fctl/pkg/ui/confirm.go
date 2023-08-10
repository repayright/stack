package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

type Confirm struct {
	style lipgloss.Style

	keyHandler *config.KeyMapHandler

	message string
}

func NewConfirm(message string) *Confirm {
	return &Confirm{
		keyHandler: config.NewKeyMapHandler(),
		style: lipgloss.NewStyle().
			Foreground(theme.SelectedColorForegroundBackground).
			Width(20).
			Height(10).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Border(lipgloss.NormalBorder()),
		message: message,
	}
}

func (c *Confirm) Styles() *lipgloss.Style {
	return &c.style
}

func (c Confirm) GetKeyMapAction() *config.KeyMapHandler {
	return c.keyHandler
}

func (c Confirm) Init() tea.Cmd {
	return nil
}

func (c Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (c Confirm) View() string {
	return c.style.Render(c.message)
}
