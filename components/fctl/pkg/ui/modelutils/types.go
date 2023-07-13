package modelutils

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model interface {
	// Tea Model
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	View() string
}
