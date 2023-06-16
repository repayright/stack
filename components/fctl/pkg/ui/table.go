package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	MaxTableHeight = 10
)

// This model implement the table.Model interface from Bubble Tea Framework
func (t TableModel) Init() tea.Cmd { return nil }

type TableModel struct {
	table table.Model
	help  bool
}

func (t TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return t, tea.Quit
		}
	}
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t TableModel) helpView() string {
	return HelpStyle.Render("Formance CLI: \n • ↑/↓: Navigate \n • q: Quit")
}

func (t TableModel) View() string {
	if t.help {
		return t.helpView() + "\n" + BaseStyle.Render(t.table.View())
	}
	return BaseStyle.Render(t.table.View()) + "\n"
}

func (t *TableModel) WithDefaultStyle() *TableModel {
	//Default styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(TabBorderColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(SelectedColorForeground).
		Background(SelectedColorForegroundBackground).
		Bold(false)

	t.table.SetStyles(s)

	return t
}

func NewDefaultOptions(ac ArrayColumn, row []table.Row) []table.Option {
	return []table.Option{
		table.WithColumns(ac),
		table.WithRows(row),
		table.WithFocused(true),
	}
}

func WithHeight(height int) table.Option {
	return table.WithHeight(height)
}

func NewTableModel(help bool, opts ...table.Option) *TableModel {
	return (&TableModel{
		table: table.New(opts...),
		help:  help,
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
