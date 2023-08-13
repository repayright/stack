package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/components/list"
	"github.com/formancehq/fctl/pkg/modelutils"
)

var (
	Regions     = "region: "
	Org         = "organization: "
	FctlVersion = "version: "
	Profile     = "profile: "
)

type Context struct {
	regions     string
	org         string
	fctlversion string
	profile     string
	model       *list.PointList
}

func NewContext() *Context {
	return &Context{
		regions:     "us-east-1",
		org:         "123456789012",
		profile:     "default",
		fctlversion: FctlVersion,
	}
}

func (c *Context) SetRegions(regions string) {
	c.regions = regions
}

func (c *Context) SetProfile(profile string) {
	c.profile = profile
}

func (c *Context) SetOrg(org string) {
	c.org = org
}

func (c *Context) SetFctlVersion(fctlversion string) {
	c.fctlversion = fctlversion
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
	Profile = modelutils.FillCharBeforeChar(Profile, " ", ":", maxWidth)

	l := []*list.HorizontalItem{
		list.NewHorizontalItem(Regions, c.regions),
		list.NewHorizontalItem(Org, c.org),
		list.NewHorizontalItem(Profile, c.profile),
		list.NewHorizontalItem(FctlVersion, c.fctlversion),
	}

	return list.NewPointList(l)
}

func (c *Context) GetMaxPossibleHeight() int {
	return 3
}

func (c *Context) GetMaxPossibleWidth() int {
	return c.model.GetMaxPossibleWidth()
}

func (c *Context) GetListKeyMapHandler() *config.KeyMapHandler {
	return nil
}

func (c *Context) Init() tea.Cmd {
	c.model = c.GeneratePointList()
	c.model.SortDir(true)
	return nil
}

func (c *Context) Update(msg tea.Msg) (*list.PointList, tea.Cmd) {
	return c.model.Update(msg)
}

func (c *Context) View() string {
	return c.model.View()
}
