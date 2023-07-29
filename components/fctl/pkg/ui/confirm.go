package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

type Confirm struct {
	keyHandler *config.KeyMapHandler
}

func NewConfirm() *Confirm {
	return &Confirm{
		keyHandler: config.NewKeyMapHandler(),
	}
}

func (c Confirm) GetKeyMapAction() *config.KeyMapHandler {
	return c.keyHandler
}
func (c Confirm) View() string {

	style := lipgloss.NewStyle().Foreground(theme.SelectedColorForegroundBackground).Width(20).Align(lipgloss.Center).Border(lipgloss.NormalBorder())

	return style.Render("Are you sure ?")
}

func (c Confirm) Init() tea.Cmd {
	return nil
}

func (c Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}
