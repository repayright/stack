package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/components/list"
	"github.com/formancehq/fctl/pkg/theme"
)

type Suggestions struct {
	list *list.PointList

	style lipgloss.Style
}

func NewSuggestions(list *list.PointList) *Suggestions {
	return &Suggestions{
		list: list,
		style: lipgloss.NewStyle().
			Foreground(theme.SelectedColorForegroundBackground).
			Align(lipgloss.Left).
			Border(lipgloss.NormalBorder()),
	}
}

func (s *Suggestions) Init() tea.Cmd {

	s.style.MaxWidth(s.list.GetMaxPossibleWidth())

	return s.list.Init()
}

func (s *Suggestions) Update(msg tea.Msg) (*Suggestions, tea.Cmd) {
	list, cmd := s.list.Update(msg)
	s.list = list
	return s, cmd
}

func (s *Suggestions) View() string {
	return s.style.Render(s.list.View())
}
