package table

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/helpers"
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
	Log := helpers.NewLogger("SCOPE")
	Log.Log(fmt.Sprintf("%d", t.Height-4))
	height := t.Height - 4
	cursorY := c.y
	if cursorY == 0 || cursorY == 1 {
		if height == len(rows) {
			Log.Log(fmt.Sprintf("Cursor ==0 && Cursor == 1:%d", len(r[0:])))
			return r[0:]
		}

		// Log.Log(fmt.Sprintf("Cursor ==0 && Cursor == 1:%d", len(r[0:height])))
		return r[0:height]
	}

	if cursorY > 1 {
		// if cursorY == len(rows) {
		// 	rows = r.Reverse()
		// 	rows = rows[cursorY-height : cursorY]
		// 	rows = rows.Reverse()
		// 	return rows
		// }

		norme := cursorY - height
		Log.Log(fmt.Sprintf("Vec1: %t", norme > 0))
		Log.Log(fmt.Sprintf("Vec2: %t", norme < 0))
		Log.Log(fmt.Sprintf("Vec1 == Vec2: %t", norme == 0))
		if norme > 0 { // Vec1: je monte
			norme = int(math.Abs(float64(norme)))
			rows = r.Reverse()
			rows = rows[height : len(rows)-cursorY]
			rows = rows.Reverse()

			return
		}

		if norme == 0 { // Vec1 == Vec2, je suis au milieu
			mid := height / 2
			rows = rows[cursorY-mid : cursorY+mid]
			// rows = rows.Reverse()
			return rows
		}

		if norme < 0 { // Vec2: je descend, c'est monté à l'envers :) ??
			rows = r.Reverse()
			if height < len(rows) {
				rows = rows[:height]
			}
			rows = rows.Reverse()
			return rows
		}

		return rows
	}

	return rows
}

func (r StyleRows) Render(c Cursor, t tea.WindowSizeMsg) string {
	out := []string{}
	// The selection start a cursorY + 1
	cursorY := c.y

	// Switch selection
	for j, row := range r.rows {
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
	}

	rows := r.rows.GetScopeRows(c, t)
	// Handle display scope
	for _, row := range rows {
		// log.Log("j", strconv.Itoa(j), "y", strconv.Itoa(t.cursor.y))
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
