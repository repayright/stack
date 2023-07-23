package ui

import (
	blist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/formancehq/fctl/pkg/ui/list"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

var (
	Regions     = "Regions: "
	Org         = "Organization: "
	FctlVersion = "Fctl Version: "
)

type Context struct {
	regions     string
	org         string
	fctlversion string
	model       *list.PointList
}

func NewContext() *Context {
	c := &Context{
		regions:     "us-east-1",
		org:         "123456789012",
		fctlversion: "0.0.1",
	}

	c.model = c.GeneratePointList()

	return c
}

func (c *Context) GeneratePointList() *list.PointList {
	maxWidth := modelutils.GetMaxCharPosXinCharList([]string{
		Regions,
		Org,
		FctlVersion,
	}, ":") + 1

	// Format to ":"
	Regions = modelutils.FillCharBeforeChar(Regions, " ", ":", maxWidth)
	Org = modelutils.FillCharBeforeChar(Org, " ", ":", maxWidth)
	FctlVersion = modelutils.FillCharBeforeChar(FctlVersion, " ", ":", maxWidth)

	return list.NewPointList(
		[]blist.Item{
			list.NewHorizontalItem(Regions, c.regions),
			list.NewHorizontalItem(Org, c.org),
			list.NewHorizontalItem(FctlVersion, c.fctlversion),
		},
		list.NewHorizontalItemDelegate(),
		maxWidth,
		4,
	)
}

func (c *Context) GetMaxPossibleHeight() int {
	return 3
}

func (c *Context) GetMaxPossibleWidth() int {
	return c.model.GetMaxPossibleWidth()
}

func (c *Context) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}

func (c *Context) Init() tea.Cmd {
	return nil
}

func (c *Context) Update(msg tea.Msg) (*list.PointList, tea.Cmd) {
	return c.model.Update(msg)
}

func (c *Context) View() string {
	return c.model.View()
}
