package table

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Cell struct {
	content string

	minWidth int

	hidden bool

	style lipgloss.Style
	Trim  *Trim
}

func NewCell(content string, opts ...CellOption) *Cell {
	c := &Cell{
		content:  content,
		minWidth: len(content),
		style:    lipgloss.NewStyle(),
	}

	for _, opt := range opts {
		c = opt(c)
	}

	return c
}

type CellOption func(*Cell) *Cell

func WithWidth(width int) CellOption {
	return func(c *Cell) *Cell {
		c.minWidth = width
		c.style = c.style.Width(width)
		c.style = c.style.MaxWidth(width)
		return c
	}
}

func (c Cell) String() string {
	return c.content
}

func (c Cell) getMinWidth() int {
	return c.minWidth
}

func (c Cell) Width() int {
	return c.style.GetWidth()
}

func (c *Cell) TrimLeft(width int) *Cell {
	c.Trim = NewTrim(width, c)
	return c
}

func (c *Cell) UnTrimLeft() *Cell {
	c.Trim = nil
	return c
}

func (c Cell) TrimRight(width int) Cell {
	c.style.MaxWidth(width)
	return c
}

func (c Cell) UnTrimRight() Cell {
	c.style.UnsetMaxWidth()
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

	return c.style.Render(c.content)
}
