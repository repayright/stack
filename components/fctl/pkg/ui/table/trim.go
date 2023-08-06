package table

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Trim struct {
	Left int
	c    *Cell
}

func NewTrim(left int, c *Cell) *Trim {
	return &Trim{
		Left: left,
		c:    c,
	}
}

func withTrim(width int, c *Cell) string {
	if width >= len(c.content) {
		dif := c.style.GetWidth() - width
		if dif < 0 {
			return ""
		}
		return strings.Repeat(" ", dif)
	}
	return c.content[width:]
}

// It should never happen because we should not render the cell
func withMaxWith(width int, c *Cell) int {
	// log := helpers.NewLogger("TRIM")
	// log.Log("width", strconv.Itoa(width), "cell", strconv.Itoa(c.style.GetWidth()+c.style.GetHorizontalMargins()+c.style.GetHorizontalPadding()))
	// log.Log("len", strconv.Itoa(len(c.content)))
	if width >= len(c.content) {
		dif := c.style.GetWidth() + c.style.GetHorizontalMargins() + c.style.GetHorizontalPadding() - width
		return dif
	}
	return c.style.GetMaxWidth() - width
}

func (t *Trim) Init() tea.Cmd {
	return nil
}

func (t *Trim) Update(msg tea.Msg) (*Trim, tea.Cmd) {
	return t, nil
}

func (t *Trim) View() string {
	return t.c.style.MaxWidth(withMaxWith(t.Left, t.c)).Render(withTrim(t.Left, t.c))
}
