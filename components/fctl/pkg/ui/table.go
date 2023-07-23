package ui

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/formancehq/fctl/pkg/ui/theme"
	"github.com/formancehq/fctl/pkg/utils"
	"golang.org/x/crypto/ssh/terminal"
)

type TableModel struct {
	columns ArrayColumn
	table   table.Model
}

// This model implement the table.Model interface from Bubble Tea Framework
func (t TableModel) Init() tea.Cmd {
	return nil
}

func (t TableModel) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	k := modelutils.NewKeyMapHandler()
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "Quit the application"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "move up"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "move down"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("? ", "Toggle help"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "show selected item"),
		),
	)

	return k
}

func (t TableModel) Update(msg tea.Msg) (Model, tea.Cmd) {
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

	minWidthBuffer := minWidthBuffer(buffer)
	Tofill := tabWidth - minWidthBuffer
	each := Tofill / len(buffer)

	for i, _ := range buffer {
		buffer[i] = buffer[i] + each
	}

	return buffer
}

func minWidthBuffer(buffer []int) int {
	count := 0
	for _, str := range buffer {
		count += str
	}
	return count
}
