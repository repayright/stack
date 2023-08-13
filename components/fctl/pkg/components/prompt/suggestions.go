package prompt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/components/table"
	"github.com/formancehq/fctl/pkg/helpers"
	"github.com/formancehq/fctl/pkg/theme"
)

type Suggestions struct {
	cursor *table.Cursor
	rows   table.StyleRows

	style lipgloss.Style
}

func NewSuggestions(list ...*table.Row) *Suggestions {
	style := table.NewStyle()
	return &Suggestions{
		cursor: table.NewCursor(),
		rows: table.NewStyleRows(
			style.Row,
			style.RowSelected,
			list...,
		),
		style: lipgloss.NewStyle().
			Foreground(theme.SelectedColorForegroundBackground).
			Align(lipgloss.Left).
			MaxHeight(4),
	}
}

func (s *Suggestions) IsEmpty() bool {
	return len(s.rows.GetRows()) == 0
}

func (s *Suggestions) Init() tea.Cmd {

	s.style.Height(len(s.rows.GetRows()))

	// for _, row := range s.rows.GetRows() {
	// 	row.
	// }

	return s.rows.Init()
}

func (s *Suggestions) Update(msg tea.Msg) (*Suggestions, tea.Cmd) {
	rows, cmd := s.rows.Update(msg)
	s.rows = rows
	return s, cmd
}

func (s *Suggestions) View() string {
	Log := helpers.NewLogger("SUGGESTIONS")
	Log.Log("View")
	Log.Log(fmt.Sprintf("%d", s.cursor.GetY()))
	return s.style.Render(s.rows.Render(s.cursor, tea.WindowSizeMsg{
		Width:  s.style.GetMaxWidth(),
		Height: s.style.GetMaxHeight(),
	}))
}
