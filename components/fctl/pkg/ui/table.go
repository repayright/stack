package ui

import (
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/theme"
	"github.com/formancehq/fctl/pkg/utils"
	"golang.org/x/crypto/ssh/terminal"
)

type TableModel struct {
	columns ArrayColumn
	table   table.Model
}

func (t TableModel) Init() tea.Cmd {
	return nil
}

func (t TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c": // Need to be deported in the key handler, same for action
			return t, tea.Quit
		}
	case tea.WindowSizeMsg:
		t.table.SetColumns(WithFullScreenTable(t.columns))
		t.table.SetHeight(msg.Height)
		return t, nil
	}

	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t TableModel) View() string {
	return theme.BaseStyle.Render(t.table.View()) + "\n"
}

func (t *TableModel) WithDefaultStyle() *TableModel {
	//Default styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.TabBorderColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(theme.SelectedColorForeground).
		Background(theme.SelectedColorForegroundBackground).
		Bold(false)

	t.table.SetStyles(s)

	return t
}

func (t *TableModel) SelectedRow() table.Row {
	return t.table.SelectedRow()
}

func WithFullScreenTable(ac ArrayColumn) ArrayColumn {

	// actual := t.table.Width()
	columnOrderedWidths := utils.Map(ac, func(c table.Column) int {
		return c.Width
	})

	// Calculate column widths
	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	margins := -15

	columnWidths := CalculateColumnWidths(columnOrderedWidths, terminalWidth+margins)
	for i := range ac {
		ac[i].Width = columnWidths[i]
	}

	return ac
}

func NewTableOptions(ac ArrayColumn, row []table.Row) []table.Option {
	return []table.Option{
		table.WithColumns(ac),
		table.WithRows(row),
		table.WithFocused(true),
	}
}

func WithHeight(height int) table.Option {
	return table.WithHeight(height)
}

func NewTableModel(columns ArrayColumn, opts ...table.Option) *TableModel {
	return (&TableModel{
		table:   table.New(opts...),
		columns: columns,
	}).WithDefaultStyle()
}

type ArrayColumn []table.Column

func NewArrayColumn(c ...table.Column) ArrayColumn {
	return c
}

func (ac ArrayColumn) AddColumn(name string, width int) ArrayColumn {
	return append(ac, NewColumn(name, width))
}
func NewColumn(name string, width int) table.Column {
	return table.Column{
		Title: name,
		Width: width,
	}
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