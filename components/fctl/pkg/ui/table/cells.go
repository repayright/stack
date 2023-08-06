package table

type Cells []*Cell

func (c Cells) AddCell(cell *Cell) Cells {
	return append(c, cell)
}

func (c Cells) Reverse() Cells {
	cells := make([]*Cell, len(c))
	for i, cell := range c {
		cells[i] = cell
	}
	return cells
}
