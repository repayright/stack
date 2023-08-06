package table

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/utils"
)

type Row struct {
	cells Cells

	activable bool
	selected  bool

	termWidth int
	style     *lipgloss.Style
}

func NewRow(items ...*Cell) *Row {
	return &Row{
		cells: items,
	}
}
func (r Row) String() string {
	mapCell := utils.Map(r.cells, func(c *Cell) string {
		return c.String()
	})
	return strings.Join(mapCell, " ")
}

func (r Row) Items() []*Cell {
	return r.cells
}

func (r Row) Init() tea.Cmd {
	return nil
}

func (r Row) Update(msg tea.Msg) (Row, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case tea.WindowSizeMsg:
		r.termWidth = msg.Width

	}
	return r, nil
}

func (r Row) View() string {
	cells := make([]string, len(r.cells))
	for i, c := range r.cells {
		cells[i] = c.View()
	}

	if r.selected {
		return r.style.Render(cells...)
	}

	return r.style.Render(cells...)
}

// func (r *Row) TrimLeft(c Cursor) *Row {
// 	width := c.x
// 	for _, cell := range r.cells {
// 		//Reset hidden
// 		cell.hidden = false

// 		// If the cell is hidden, skip it
// 		if cell.Width() > w {
// 			w -= cell.Width()
// 			cell.hidden = true
// 			continue
// 		}

// 		cell.TrimLeft(w)
// 	}

// 	return r
// }

// func (r *Row) TrimRight(c Cursor) *Row {
// 	width := c.x
// 	for _, cell := range r.cells {
// 		//Reset hidden
// 		cell.hidden = false

// 		// If the cell is hidden, skip it
// 		if cell.Width() > w {
// 			w -= cell.Width()
// 			cell.hidden = true
// 			continue
// 		}

// 		cell.TrimRight(w)
// 	}

// 	return r
// }

// This function need to decide if we need to trim left or right
// Depending on the cursor position
// And the row width
func (r *Row) Trim(c Cursor) *Row {
	x := c.x
	for i, cell := range r.cells {
		if x == 0 { // We would need to set a max width to termsize

			break
		}

		cellWidth := cell.Width()
		if i*x >= cellWidth {
			cell.hidden = true
			x -= cellWidth
		}
	}

	return r
}
