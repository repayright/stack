package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/components/table"
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
			MaxHeight(len(list)),
	}
}

func (s *Suggestions) IsEmpty() bool {
	return len(s.rows.GetRows()) == 0
}

func (s *Suggestions) Init() tea.Cmd {

	s.style.Height(len(s.rows.GetRows())).MaxHeight(len(s.rows.GetRows()))
	return s.rows.Init()
}

func (s Suggestions) GetSelected() *table.Row {
	return s.rows.GetRows().GetSelected(s.cursor)
}

func (s Suggestions) Update(msg tea.Msg) (Suggestions, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if cursorY := s.cursor.GetY(); cursorY-1 >= 0 {
				s.cursor.MoveUp()
			}
		case "down", "j", "tab":
			if cursorY := s.cursor.GetY(); cursorY+1 == len(s.rows.GetRows()) {
				s.cursor.ToTop()
				return s, nil
			}
			if cursorY := s.cursor.GetY(); cursorY+1 < len(s.rows.GetRows()) {
				s.cursor.MoveDown()
				return s, nil
			}
		}
	default:
		rows, cmd := s.rows.Update(msg)
		s.rows = rows
		return s, cmd
	}

	return s, nil
}

func (s Suggestions) View() string {
	// Log := helpers.NewLogger("SUGGESTIONS")
	// Log.Log("View")
	// Log.Log(fmt.Sprintf("%d", s.cursor.GetY()))
	return s.style.Render(s.rows.Render(s.cursor, tea.WindowSizeMsg{
		Width:  s.style.GetMaxWidth(),
		Height: s.style.GetMaxHeight(),
	}))
}
