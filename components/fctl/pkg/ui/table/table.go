package table

import (
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/formancehq/fctl/pkg/utils"
)

type TableOption func(*Table) *Table

type Table struct {
	// Column
	options []TableOption

	// Row
	header *Row
	rows   Rows

	style *Style

	// For full screen table
	fullScreen   bool
	minWidth     []int
	terminalSize tea.WindowSizeMsg

	// Column Cursor
	cursor *Cursor
}

func (t Table) renderRow(r *Row) string {
	// log := helpers.NewLogger("ROW")
	buffer := t.cursor.x
	for _, c := range r.cells {
		if buffer < 0 {
			//Just want to break the first for loop in a magic case
			break
		}

		width := c.style.GetMaxWidth() + c.style.GetHorizontalPadding() + c.style.GetHorizontalMargins()
		// log.Log("i", strconv.Itoa(i), "width : ", strconv.Itoa(width), ", x", strconv.Itoa(t.cursor.x))

		//Reset hidden cells
		c.hidden = false
		c.UnTrimLeft()

		if buffer >= width {
			c.hidden = true
			buffer -= width
			continue
		}

		if buffer < width && buffer >= 0 {
			// log.Log("i", strconv.Itoa(i), "buffer", strconv.Itoa(buffer))
			c.TrimLeft(buffer)
			buffer -= width
		}

		// //Allign everything to header style cells
		// r.cells[i].style = r.cells[i].style.Width(c.Width())
	}

	return r.View()
}

func (t Table) renderRows() string {
	out := []string{}
	// log := helpers.NewLogger("ROWS")
	rows := t.rows[:]
	for j, r := range rows {
		// log.Log("j", strconv.Itoa(j), "y", strconv.Itoa(t.cursor.y))

		// if y < t.cursor.y {
		// 	continue
		// }

		// (i,j) is for the cursor selection,
		// (j+1) is for the header selection,
		// we should consider the header as a row in order to be able to sort columns
		if j+1 == t.cursor.y {
			r.style = &t.style.RowSelected
		} else {
			r.style = &t.style.Row
		}

		// // We should consider the cursor here
		// // and hidden columns wich mean hidden cells
		// // We should use TrimLeft and TrimRight
		// if t.cursor.x > len(r.cells)-1 {
		// 	continue
		// }

		style := r.style.MaxWidth(t.terminalSize.Width - 3)
		r.style = &style
		out = append(out, t.renderRow(r))
	}
	return lipgloss.JoinVertical(lipgloss.Top, out...)
}
func (t Table) renderHeader() string {
	// t.header.TrimLeft(t.cursor.x)
	var style lipgloss.Style
	if t.cursor.y == 0 {
		style = t.style.HeaderSelected
	} else {
		style = t.style.Header.Background(lipgloss.Color("#1D1D1D"))
	}
	style = style.MaxWidth(t.terminalSize.Width - 4)

	t.header.style = &style

	return t.renderRow(t.header)
}
func (t *Table) Init() tea.Cmd {

	var cmd tea.Cmd

	for _, opt := range t.options {
		t = opt(t)
	}

	t.header.style = &t.style.Header
	cmd = t.header.Init()

	// This is a hack to fix the header width
	// TODO: Find a better way to do this
	// Shit happen between Header Cells n
	for _, c := range t.header.cells {
		c.content += " "
		// c.style.MarginRight(1)
	}

	w, h, err := modelutils.GetTerminalSize()
	if err != nil {
		panic(err)
	}

	t.terminalSize = tea.WindowSizeMsg{Width: w, Height: h}

	cmd = tea.Batch(cmd, t.rows.Init(), func() tea.Msg {
		return t.terminalSize
	})

	for _, r := range t.rows {
		for i, c := range r.cells {
			WithWidth(t.header.cells[i].getMinWidth() + t.header.cells[0].style.GetHorizontalMargins())(c)
		}
	}

	return cmd
}

func (t Table) ResetCursor() {
	t.cursor.x = 0
	t.cursor.y = 0
}

func (t Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.terminalSize = msg

		return &t, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			t.cursor.x = utils.Max(t.cursor.x-1, 0)
		case "right":
			t.cursor.x = utils.Min(t.cursor.x+1, t.terminalSize.Width-3)
		case "down":
			t.cursor.y = utils.Min(t.cursor.y+1, len(t.rows))
		case "up":
			t.cursor.y = utils.Max(t.cursor.y-1, 0)
		}
	case changeScopeMsg:
		return &t, func() tea.Msg {
			return t.terminalSize
		}
	}

	t.rows, cmd = t.rows.Update(msg)

	return &t, tea.Batch(cmd)
}

func (t Table) View() string {
	render :=
		[]string{
			t.renderHeader(),
			t.renderRows(),
		}

	border := t.style.Wrapper.Width(t.terminalSize.Width - 2)
	innerBox := t.style.Body.MaxWidth(t.terminalSize.Width - 3)

	return border.Render(lipgloss.PlaceHorizontal(innerBox.GetWidth(), 0, lipgloss.JoinVertical(lipgloss.Top, render...)))
}

func getBoolDir(b bool) int {
	if b {
		return 1
	}
	return -1
}

func (t *Table) SelectedRow() *Row {

	return t.rows[t.cursor.y]
}

func WithDefaultStyle() TableOption {
	return func(t *Table) *Table {
		t.style = NewStyle()
		for i := range t.rows {
			t.rows[i].style = &t.style.Row
		}

		t.header.style = &t.style.Header
		return t
	}
}

func NewTable(header *Row, rows []*Row, tableOptions ...TableOption) *Table {
	t := &Table{
		header:  header,
		rows:    rows,
		options: tableOptions,
		cursor:  NewCursor(),
	}

	for _, opt := range tableOptions {
		table := opt(t)
		t = table
	}

	return t
}

func CalculateColumnWidths(buffer []int, tabWidth int) []int {

	minWidthBuffer := minWidthIntList(buffer)
	Tofill := tabWidth - minWidthBuffer
	each := Tofill / len(buffer)
	for i := range buffer {
		buffer[i] = buffer[i] + each
	}

	return buffer
}

func minWidthIntList(list []int) int {
	count := 0
	for _, str := range list {
		count += str
	}
	return count
}
