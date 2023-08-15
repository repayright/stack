package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/components/table"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/modelutils"
	"github.com/formancehq/fctl/pkg/theme"
)

type Confirm struct {
	cursor *table.Cursor

	style             lipgloss.Style
	cellStyle         lipgloss.Style
	cellStyleSelected lipgloss.Style
	row               *table.Row

	keyHandler *config.KeyMapHandler

	message string
	action  tea.Cmd
}

func NewConfirm(msg modelutils.ConfirmActionMsg) *Confirm {
	width := 40

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Align(lipgloss.Center)
		// Background(theme.SelectedColorForegroundBackground)
	cellStyleSelected := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(theme.SelectedColorForegroundBackground)

	rowStyle := lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Align(lipgloss.Center).
		Background(theme.GetBackgroundColor())

	cells := table.NewCells(
		table.NewCell("No", table.WithStyle(cellStyleSelected)),
		table.NewCell("Yes", table.WithStyle(cellStyle)),
	)

	return &Confirm{
		cursor: table.NewCursor(
			table.WithX(0),
		),
		style: lipgloss.NewStyle().
			Foreground(theme.SelectedColorForegroundBackground).
			Width(width).
			Height(10).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Foreground(theme.BlackColor).
			Border(lipgloss.NormalBorder()),
		cellStyle:         cellStyle,
		cellStyleSelected: cellStyleSelected,
		row: table.NewRow(
			cells,
			table.WithRowStyle(rowStyle),
		),
		message: msg.Question,
		action:  msg.Action,
	}
}

func (c *Confirm) Styles() lipgloss.Style {
	return c.style
}

func (c Confirm) GetKeyMapAction() *config.KeyMapHandler {
	return c.keyHandler
}

func (c Confirm) Init() tea.Cmd {
	return nil
}

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if c.cursor.GetX() > 0 {
				cells := c.row.Items()
				table.WithStyle(c.cellStyle)(cells[c.cursor.GetX()])
				c.cursor.MoveLeft()
				table.WithStyle(c.cellStyleSelected)(cells[c.cursor.GetX()])
			}
		case "right":
			if c.cursor.GetX() < 1 {
				cells := c.row.Items()
				table.WithStyle(c.cellStyle)(cells[c.cursor.GetX()])
				c.cursor.MoveRight()
				table.WithStyle(c.cellStyleSelected)(cells[c.cursor.GetX()])
			}
		case "enter":
			if c.cursor.GetX() == 0 {
				return c, func() tea.Msg {
					return modelutils.CloseConfirmMsg{}
				}
			}
			return c, tea.Sequence(func() tea.Msg {
				return modelutils.CloseConfirmMsg{}
			}, c.action)
		case "y":
			return c, tea.Sequence(func() tea.Msg {
				return modelutils.CloseConfirmMsg{}
			}, c.action)
		case "n":
			return c, func() tea.Msg {
				return modelutils.CloseConfirmMsg{}
			}
		}
	}

	return c, nil
}

func (c Confirm) View() string {
	return c.style.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.NewStyle().Margin(2, 2, 2, 2).Render(c.message),
			c.row.View(),
		),
	)
}
