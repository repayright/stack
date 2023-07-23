package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type Model interface {
	GetListKeyMapHandler() *modelutils.KeyMapHandler

	// Tea Model
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	View() string
}
