package table

type Cells []*Cell

func (c Cells) AddCell(cell *Cell) Cells {
	return append(c, cell)
}
