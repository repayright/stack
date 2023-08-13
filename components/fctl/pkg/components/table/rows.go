package table

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/helpers"
)

type Rows []*Row

func (r Rows) Reverse() Rows {
	rows := make([]*Row, len(r))
	for i, row := range r {
		rows[i] = row
	}
	return rows
}
func (r Rows) GetScopeY(c *Cursor, height int) (rows Rows) {

	cursorY := c.y
	if cursorY <= 0 {
		return r[0:int(math.Min(float64(len(r)), float64(height)))]
	}
	mid := height / 2
	// Return lasts
	if cursorY+mid >= len(r) {
		max := int(math.Min(float64(len(r)), float64(cursorY+mid+height%2)))
		r = r[int(math.Max(float64(len(r)-height), 0)):max]
		return r
	}

	// Return firsts
	if cursorY-mid <= 0 {
		return r[0:int(math.Min(float64(len(r)), float64(height)))]
	}

	// Return others
	// cursorY+mid+1 <=> mid + 1, it fix an int overflow about mid division
	// Thougth if 5/2 = 2.5
	// For, a = b × q + r
	// a = 5
	// q = 2
	// r = 1
	// r <=> 5%2 <=> 1
	return r[cursorY-mid : int(math.Min(float64(cursorY+mid+height%2), float64(len(r))))]

}

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
func (st StyleRows) GetRows() Rows {
	return st.rows
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

func (r StyleRows) Render(c *Cursor, t tea.WindowSizeMsg) string {
	out := []string{}
	cursorY := c.y
	Log := helpers.NewLogger("ROWS")
	Log.Log(fmt.Sprintf("%d", cursorY))

	// Switch selection
	for j, row := range r.rows {
		if cursorY == j {
			row.style = r.selectedStyle
		} else {
			row.style = r.rawStyle
		}
	}
	// Render only scoped rows
	rows := r.rows.GetScopeY(c, t.Height)

	for _, row := range rows {
		// style := row.style.Width(t.Width - row.style.GetHorizontalMargins() - row.style.GetHorizontalPadding())
		// rows[i].style = style
		out = append(out, row.Render(c))
	}
	return lipgloss.JoinVertical(lipgloss.Top, out...)
}

func (r StyleRows) View() string {
	out := []string{}
	for _, row := range r.rows {
		out = append(out, row.View())
	}
	return lipgloss.JoinVertical(lipgloss.Top, out...)
}
