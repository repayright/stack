package table

type Trim struct {
	Left int
}

func NewTrim(left int) *Trim {
	return &Trim{
		Left: left,
	}
}

func withTrim(width int, c Cell) string {
	if width > c.minWidth {
		return ""
	}
	return c.content[width:]
}

// It should never happen because we should not render the cell
func withMaxWith(width int, c Cell) int {
	if c.style.GetMaxWidth()-width < 0 {
		return 0
	}
	return c.style.GetMaxWidth() - width
}

func (t Trim) View(c Cell) string {
	return c.style.MaxWidth(withMaxWith(t.Left, c)).Render(withTrim(t.Left, c))
}
