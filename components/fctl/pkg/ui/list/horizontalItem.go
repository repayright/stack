package list

type HorizontalItem struct {
	title, desc string
}

func NewHorizontalItem(title, desc string) *HorizontalItem {
	return &HorizontalItem{
		title: title,
		desc:  desc,
	}
}

func (i HorizontalItem) GetWidth() int {
	return len(i.title) + len(i.desc) + 2 // +2 x " " between title and desc

}

func (i *HorizontalItem) GetHeight() int {
	return 1
}
