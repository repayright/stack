package table

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/modelutils"
	"github.com/formancehq/fctl/pkg/utils"
)

type Row struct {
	cells Cells

	activable bool
	hidden    bool

	termWidth int
	style     lipgloss.Style
}

type RowOpts func(*Row) *Row

func WithRowStyle(style lipgloss.Style) RowOpts {
	return func(r *Row) *Row {
		r.style = style
		return r
	}
}

func NewRow(items []*Cell, opts ...RowOpts) *Row {
	r := &Row{
		cells: items,
	}

	for _, opt := range opts {
		r = opt(r)
	}

	return r
}

func (r *Row) GetWidth() int {
	return r.style.GetWidth()
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
	var cmd tea.Cmd
	for _, cell := range r.cells {
		cmd = tea.Batch(cell.Init())
	}
	return cmd
}

func (r Row) Update(msg tea.Msg) (Row, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case tea.WindowSizeMsg:
		r.termWidth = msg.Width
	case modelutils.ResizeMsg:
		r.termWidth = msg.Width
	}
	return r, nil
}

func (r Row) View() string {
	cells := make([]string, len(r.cells))
	for i, c := range r.cells {
		if c.hidden {
			continue
		}
		cells[i] = c.View()
	}

	return r.style.Render(cells...)
}

// This function need to decide if we need to trim left or right
// Depending on the cursor position
// And the row width
func (r *Row) Trim(c Cursor) *Row {
	x := c.x
	for i, cell := range r.cells {
		cellWidth := cell.Width()
		if i*x >= cellWidth {
			cell.hidden = true
			x -= cellWidth
		}
	}

	return r
}

func (r Row) Render(c *Cursor) string {
	buffer := c.x
	for _, c := range r.cells {
		width := c.style.GetMaxWidth() - c.style.GetHorizontalPadding() - c.style.GetHorizontalMargins()

		//Reset hidden cells
		c.hidden = false
		c.UnTrimLeft()

		if buffer >= width {
			c.hidden = true
			buffer -= width
			continue
		}

		if buffer < width && buffer >= 0 {
			c.TrimLeft(buffer)
			buffer -= width
		}

	}

	return r.View()
}
func (r *Row) AddCell(cell *Cell) *Row {
	r.cells = append(r.cells, cell)
	return r
}
