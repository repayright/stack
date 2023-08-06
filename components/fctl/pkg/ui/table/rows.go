package table

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (r *Row) AddCell(cell *Cell) *Row {
	r.cells = append(r.cells, cell)
	return r
}

type Rows []*Row

type StyleRows struct {
	rows          Rows
	rawStyle      lipgloss.Style
	selectedStyle lipgloss.Style
}

func NewStyleRows(rawStyle, selectedStyle lipgloss.Style, rows ...*Row) StyleRows {
	return StyleRows{
		rows:          Rows(rows),
		rawStyle:      rawStyle,
		selectedStyle: selectedStyle,
	}
}

func (r StyleRows) AddRow(row *Row) StyleRows {
	r.rows = append(r.rows, row)
	return r
}

func (r StyleRows) Init() tea.Cmd {
	var cmd tea.Cmd
	for _, row := range r.rows {
		cmd = tea.Batch(row.Init())
	}
	return cmd
}

func (r StyleRows) Update(msg tea.Msg) (StyleRows, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	for i, row := range r.rows {
		rowi, cmd := row.Update(msg)
		r.rows[i] = &rowi
		cmds = append(cmds, cmd)
	}
	return r, tea.Batch(cmds...)
}

func (r StyleRows) Render(c Cursor, t tea.WindowSizeMsg) string {
	out := []string{}
	// log := helpers.NewLogger("ROWS")
	var rows Rows

	cursorRows := len(r.rows[c.y:])
	if cursorRows >= t.Height-3 {
		rows = r.rows[c.y:]
	} else if cursorRows < t.Height-3 {
		origin := t.Height - c.y - 3
		if origin < 0 {
			origin = 0
		}

		rows = r.rows[origin:]
	}

	for j, row := range rows {
		// log.Log("j", strconv.Itoa(j), "y", strconv.Itoa(t.cursor.y))

		if j+1 > t.Height-3 {
			row.hidden = true
			continue
		} else {

			row.hidden = false
		}

		// (i,j) is for the cursor selection,
		// (j+1) is for the header selection,
		// we should consider the header as a row in order to be able to sort columns
		if j+1 == c.y {
			row.style = r.selectedStyle
		} else {
			row.style = r.rawStyle
		}

		// // We should consider the cursor here
		// // and hidden columns wich mean hidden cells
		// // We should use TrimLeft and TrimRight
		// if t.cursor.x > len(r.cells)-1 {
		// 	continue
		// }

		style := row.style.MaxWidth(t.Width - 3)
		row.style = style
		out = append(out, row.Render(c))
	}
	return lipgloss.JoinVertical(lipgloss.Top, out...)
}
