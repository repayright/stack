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

func (r Rows) Reverse() Rows {
	rows := make([]*Row, len(r))
	for i, row := range r {
		rows[i] = row
	}
	return rows
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

func (r Rows) GetScopeRows(c Cursor, t tea.WindowSizeMsg) (rows Rows) {

	cursorY := c.y
	var cursorRows []*Row
	if cursorY == len(r) {
		cursorRows = r[t.Height-cursorY-1:]
	} else if cursorY >= 1 {
		cursorRows = r[cursorY-1:]
	} else {
		cursorRows = r[cursorY:]
	}

	if len(cursorRows) > t.Height-1 {
		rows = cursorRows
	} else if len(cursorRows) <= t.Height-1 {

		if t.Height-cursorY+1 > len(r) {
			cursorRows = r.Reverse()
			cursorRows = r[t.Height-cursorY+1:]
			cursorRows = r.Reverse()
		} else if t.Height-cursorY+1 <= len(r) {
			cursorRows = r.Reverse()
			cursorRows = r[cursorY+1:]
			cursorRows = r.Reverse()
		} else {
			cursorRows = r.Reverse()
			cursorRows = r[cursorY:]
			cursorRows = r.Reverse()
		}

		rows = cursorRows
	}
	return rows
}

func (r StyleRows) Render(c Cursor, t tea.WindowSizeMsg) string {
	out := []string{}
	// The selection start a cursorY + 1
	cursorY := c.y
	rows := r.rows.GetScopeRows(c, t)

	// Displaying
	for j, row := range rows {
		// log.Log("j", strconv.Itoa(j), "y", strconv.Itoa(t.cursor.y))

		if j+1 > t.Height-1 {
			row.hidden = true
			continue
		} else {
			row.hidden = false
		}

		// (i,j) is for the cursor selection,
		// (j+1) is for the header selection,
		// we should consider the header as a row in order to be able to sort columns

		if cursorY == 0 {
			row.style = r.rawStyle
		} else if j+1 == cursorY {
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
