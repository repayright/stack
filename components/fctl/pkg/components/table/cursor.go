package table

type Cursor struct {
	x int
	y int
}
type CursorOption func(*Cursor) *Cursor

func NewCursor(opts ...CursorOption) *Cursor {
	c := &Cursor{
		x: 0,
		y: 0,
	}

	for _, opt := range opts {
		c = opt(c)
	}

	return c
}

func WithX(x int) CursorOption {
	return func(c *Cursor) *Cursor {
		c.x = x
		return c
	}
}

func WithY(y int) CursorOption {
	return func(c *Cursor) *Cursor {
		c.y = y
		return c
	}
}

func (c Cursor) GetX() int {
	return c.x
}

func (c Cursor) GetY() int {
	return c.y
}

func (c *Cursor) MoveLeft() {
	c.x--
}

func (c *Cursor) MoveRight() {
	c.x++
}

func (c *Cursor) MoveUp() {
	c.y--
}

func (c *Cursor) MoveDown() {
	c.y++
}

func (c *Cursor) MoveTo(x, y int) {
	c.x = x
	c.y = y
}

func (c *Cursor) ToTop() {
	c.y = 0
}
