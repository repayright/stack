package renderer

import tea "github.com/charmbracelet/bubbletea"

type Rule struct {
}

type Renderer struct {
	rendered string
	models   []tea.Model
	// action   map[func() bool]tea.Model
}
