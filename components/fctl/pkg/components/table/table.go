package table

import (
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/modelutils"
	"github.com/formancehq/fctl/pkg/utils"
)

type TableOption func(*Table) *Table

type Table struct {
	// Column
	options []TableOption

	// Row
	header *Row
	rows   StyleRows

	style *Style

	// For full screen table
	fullScreen   bool
	terminalSize tea.WindowSizeMsg

	// Column Cursor
	cursor *Cursor
}

func (t Table) renderRow(r *Row) string {
	return r.Render(*t.cursor)
}

func (t Table) renderRows() string {
	return t.rows.Render(*t.cursor, tea.WindowSizeMsg{
		Width:  t.style.Body.GetHorizontalFrameSize(),
		Height: t.terminalSize.Height,
	})
}

func (t Table) renderHeader() string {
	var style lipgloss.Style
	if t.cursor.y == 0 {
		style = t.style.HeaderSelected
	} else {
		style = t.style.Header
	}

	// // -3 is for the border, margins and padding
	// style = style.MaxWidth(t.terminalSize.Width - style.GetHorizontalMargins() - style.GetHorizontalPadding() - 1)

	t.header.style = style

	return t.renderRow(t.header)
}
func (t *Table) Init() tea.Cmd {

	var cmd tea.Cmd

	for _, opt := range t.options {
		t = opt(t)
	}

	t.header.style = t.style.Header
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

	// When table is a root model
	if t.fullScreen {
		t.terminalSize = tea.WindowSizeMsg{Width: w, Height: h}
	} else {
		// Rows + headerRow
		// Default heigth is equal to rows + header
		t.terminalSize = tea.WindowSizeMsg{Width: t.header.style.GetMaxWidth(), Height: len(t.rows.rows) + 1}
	}

	cmd = tea.Batch(cmd, t.rows.Init(), func() tea.Msg {
		return t.terminalSize
	})

	// Set default with for each rows
	for _, r := range t.rows.rows {
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
			t.cursor.y = utils.Min(t.cursor.y+1, len(t.rows.rows))
		case "up":
			t.cursor.y = utils.Max(t.cursor.y-1, 0)
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

	border := t.style.Wrapper.Width(t.terminalSize.Width - t.style.Wrapper.GetHorizontalFrameSize()).Height(t.terminalSize.Height)
	innerBox := t.style.Body.MaxWidth(t.terminalSize.Width - border.GetHorizontalFrameSize()).MaxHeight(t.terminalSize.Height - border.GetVerticalFrameSize())
	content := innerBox.Render(lipgloss.PlaceHorizontal(innerBox.GetWidth(), 0, lipgloss.JoinVertical(lipgloss.Top, render...)))
	return border.Render(content)
}

func (t *Table) SelectedRow() *Row {
	if t.cursor.y == 0 {
		return t.header
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

func WithFullScreen(fullScreen bool) TableOption {
	return func(t *Table) *Table {
		t.fullScreen = fullScreen
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
