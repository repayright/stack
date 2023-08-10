package table

import (
	"math"

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

func (r Rows) GetScopeRows(c Cursor, height int) (rows Rows) {
	cursorY := c.y
	if cursorY <= 0 {
		return r[0:int(math.Min(float64(len(r)), float64(height)))]
	}

	mid := height / 2
	if cursorY+mid >= len(r) { // Ne renvoi que les derniers
		max := int(math.Min(float64(len(r)), float64(cursorY+mid)))
		r = r[int(math.Max(float64(len(r)-height), 0)):max]
		return r
	}
	if cursorY-mid <= 0 { // Ne renvoi que les premiers
		return r[0:int(math.Min(float64(len(r)), float64(height)))]
	}

	// Centre
	return r[cursorY-mid : int(math.Min(float64(cursorY+mid+1), float64(len(r))))]

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

	rows := r.rows.GetScopeRows(c, t.Height)
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
