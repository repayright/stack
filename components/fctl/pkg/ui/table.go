package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

var (
	MaxTableHeight = 10
)

// This model implement the table.Model interface from Bubble Tea Framework
func (t TableModel) Init() tea.Cmd { return nil }

type TableModel struct {
	table table.Model
}

func (t TableModel) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	k := modelutils.NewKeyMapHandler()
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q  ", "Quit the application"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
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
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)

	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)

	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
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
	}
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t TableModel) View() string {
	return fctl.BaseStyle.Render(t.table.View()) + "\n"
}

func (t *TableModel) WithDefaultStyle() *TableModel {
	//Default styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(fctl.TabBorderColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(fctl.SelectedColorForeground).
		Background(fctl.SelectedColorForegroundBackground).
		Bold(false)

	t.table.SetStyles(s)

	return t
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

func NewTableModel(opts ...table.Option) *TableModel {
	return (&TableModel{
		table: table.New(opts...),
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
