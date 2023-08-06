package table

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Row) AddCell(cell *Cell) *Row {
	r.cells = append(r.cells, cell)
	return r
}

type Rows []*Row

func (r Rows) AddRow(row *Row) Rows {
	return append(r, row)
}

func (r Rows) Init() tea.Cmd {
	var cmd tea.Cmd
	for _, row := range r {
		cmd = tea.Batch(row.Init())
	}
	return cmd
}

func (rows Rows) Update(msg tea.Msg) (Rows, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	for _, row := range rows {
		_, cmd := row.Update(msg)
		cmds = append(cmds, cmd)
	}
	return rows, tea.Batch(cmds...)
}
