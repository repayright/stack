package table

import (
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/utils"
)

type TableOption func(*Table) *Table

type Table struct {
	// Column
	options []TableOption

	// Should never mutate this value
	minHeaderWidth int

	// Row
	header *Row
	rows   StyleRows

	style *Style

	terminalSize tea.WindowSizeMsg

	// Column Cursor
	cursor *Cursor
}

func (t Table) renderRow(r *Row) string {
	return r.Render(t.cursor)
}

func (t Table) renderRows() string {
	// Cursor for rows table start at 1
	cursor := NewCursor(WithY(t.cursor.y-1), WithX(t.cursor.x))

	return t.rows.Render(cursor, tea.WindowSizeMsg{
		Width:  t.style.Body.GetMaxWidth(),
		Height: utils.Max(t.style.Body.GetMaxHeight()-1, 0),
	})
}

func (t Table) renderHeader() string {
	var style lipgloss.Style
	if t.cursor.y == 0 {
		style = t.style.HeaderSelected
	} else {
		style = t.style.Header
	}

	t.header.style = style.MaxWidth(t.terminalSize.Width - t.style.Wrapper.GetHorizontalFrameSize() - 1)

	return t.renderRow(t.header)
}
func (t *Table) Init() tea.Cmd {

	var cmd tea.Cmd

	for _, opt := range t.options {
		t = opt(t)
	}

	t.header.style = t.style.Header
	cmd = t.header.Init()
	t.minHeaderWidth = t.getHeaderSize()
	for _, r := range t.rows.rows {
		for i, c := range r.cells {
			WithWidth(t.header.cells[i].getMinWidth() + t.header.cells[0].style.GetHorizontalMargins())(c)
		}
	}

	return cmd
}

func (t Table) getHeaderSize() int {
	size := 0
	for _, cell := range t.header.Items() {
		size += cell.Width()
	}

	return size
}

func (t Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.cursor.x = 0
		t.terminalSize = msg
		t.style.Wrapper = t.style.Wrapper.Width(t.terminalSize.Width - t.style.Wrapper.GetHorizontalFrameSize()).
			Height(t.terminalSize.Height)
		t.style.Body = t.style.Body.MaxWidth(t.terminalSize.Width - t.style.Wrapper.GetHorizontalFrameSize()).
			MaxHeight(t.style.Wrapper.GetHeight())

		actualHeaderSize := t.getHeaderSize()
		if actualHeaderSize > t.style.Body.GetMaxWidth() && actualHeaderSize >= t.minHeaderWidth {
			for _, c := range t.header.cells {
				c.SetFill(0)
				c.UnTrimLeft()
			}
			for _, r := range t.rows.rows {
				for _, c := range r.cells {
					c.SetFill(0)
					c.UnTrimLeft()
				}
			}
		}

		actualHeaderSize = t.getHeaderSize()
		if actualHeaderSize < t.style.Body.GetMaxWidth() &&
			actualHeaderSize >= t.minHeaderWidth {

			bodyWidth := t.style.Body.GetMaxWidth()
			toFill := bodyWidth - actualHeaderSize

			for i := 0; i < toFill; i++ {
				c := t.header.cells[i%len(t.header.cells)]
				c.SetFill(c.GetFill() + 1)
				t.header.cells[i%len(t.header.cells)] = c
			}
			for _, r := range t.rows.rows {
				for i, c := range r.cells {
					c.SetFill(t.header.cells[i].GetFill())
				}
			}
		}

		return &t, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			t.cursor.x = utils.Max(t.cursor.x-2, 0)
		case "right":
			t.cursor.x = utils.Min(t.cursor.x+2, t.getHeaderSize()-t.style.Body.GetMaxWidth()-t.style.Body.GetHorizontalFrameSize()+t.header.style.GetHorizontalFrameSize()+4) //
		case "down":
			t.cursor.y = utils.Min(t.cursor.y+1, len(t.rows.rows))
		case "up":
			t.cursor.y = utils.Max(t.cursor.y-1, 0)
		case "tab":
			t.cursor.y = t.cursor.y + 1
			if t.cursor.y == len(t.rows.rows)+1 {
				t.cursor.y = 0
			}
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

	content := t.style.Body.
		Render(lipgloss.PlaceHorizontal(t.style.Body.GetWidth(),
			0,
			lipgloss.JoinVertical(lipgloss.Top, render...),
		))

	return t.style.Wrapper.Render(content)
}

func (t *Table) SelectedRow() *Row {
	// For the moment we don't support header selection
	if t.cursor.y == 0 {
		return nil
	}

	if t.cursor.y == 1 {
		return t.rows.rows[0]
	}

	return t.rows.rows[t.cursor.y-1]
}

func WithDefaultStyle() TableOption {
	return func(t *Table) *Table {
		t.style = NewStyle()
		for i := range t.rows.rows {
			t.rows.rows[i].style = t.style.Row
		}

		t.header.style = t.style.Header
		return t
	}
}

func WithSize(msg tea.WindowSizeMsg) TableOption {
	return func(t *Table) *Table {
		t.terminalSize = msg
		return t
	}
}

func NewTable(header *Row, rows []*Row, tableOptions ...TableOption) *Table {
	t := &Table{
		header: header,
		rows: NewStyleRows(
			NewStyle().Row,
			NewStyle().RowSelected,
			rows...,
		),
		options: tableOptions,
		cursor:  NewCursor(),
	}

	for _, opt := range tableOptions {
		table := opt(t)
		t = table
	}

	return t
}
