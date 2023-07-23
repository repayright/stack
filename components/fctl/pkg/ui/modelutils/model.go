package modelutils

import tea "github.com/charmbracelet/bubbletea"

// Model tea.Model interface
type Model interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}
