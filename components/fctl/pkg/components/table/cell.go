package table

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/utils"
)

type Cell struct {
	content string

	minWidth int
	fill     int

	hidden bool

	style lipgloss.Style
	Trim  *Trim
}

func NewCell(content string, opts ...CellOption) *Cell {
	c := &Cell{
		content:  content,
		minWidth: len(content),
		style:    lipgloss.NewStyle(),
		fill:     0,
	}

	for _, opt := range opts {
		c = opt(c)
	}

	return c
}

type CellOption func(*Cell) *Cell

func WithStyle(style lipgloss.Style) CellOption {
	return func(c *Cell) *Cell {
		c.style = style
		return c
	}
}

func WithWidth(width int) CellOption {
	return func(c *Cell) *Cell {
		c.minWidth = width
		c.style = c.style.Width(width + c.fill)
		c.style = c.style.MaxWidth(width + c.fill)
		return c
	}
}
func (c *Cell) Content() string {
	return c.content
}

func (c *Cell) String() string {
	return c.content + strings.Repeat(" ", c.fill)
}

func (c *Cell) getMinWidth() int {
	return c.minWidth
}

func (c *Cell) Width() int {
	return c.style.GetWidth()
}

func (c *Cell) GetFill() int {
	return c.fill
}

func (c *Cell) SetFill(fill int) *Cell {
	fill = utils.Max(fill, 0)

	c.fill = fill
	c.style = c.style.Width(c.minWidth + fill)
	c.style = c.style.MaxWidth(c.minWidth + fill)

	return c
}

func (c *Cell) TrimLeft(width int) *Cell {
	c.Trim = NewTrim(width, c)
	return c
}

func (c *Cell) UnTrimLeft() *Cell {
	c.Trim = nil
	return c
}

func (c Cell) Init() tea.Cmd {
	return nil
}

func (c Cell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd tea.Cmd
	)
	if c.Trim != nil {
		_, cmd = c.Trim.Update(msg)
	}

	return nil, cmd
}

func (c Cell) View() string {
	if c.hidden {
		return ""
	}
	if c.Trim != nil {
		return c.Trim.View()
	}

	return c.style.Render(c.String())
}
