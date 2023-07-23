package profiles

import (
	"flag"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList         = "list"
	descriptionList = "List all profiles"
)

type Profile struct {
	Name   string `json:"name"`
	Active string `json:"active"`
}
type ListStore struct {
	Profiles []*Profile `json:"profiles"`
}

func NewListStore() *ListStore {
	return &ListStore{
		Profiles: []*Profile{},
	}
}

func NewListConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	return config.NewControllerConfig(
		useList,
		descriptionList,
		descriptionList,
		[]string{
			"ls",
			"l",
		},
		flags,
	)
}

var _ config.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *config.ControllerConfig
}

func NewListController(config *config.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	profiles := cfg.GetProfiles()

	p := fctl.MapKeys(profiles)
	currentProfileName := fctl.GetCurrentProfileName(flags, cfg)

	for _, k := range p {
		c.store.Profiles = append(c.store.Profiles, &Profile{
			Name: k,
			Active: func(k string) string {
				if currentProfileName == k {
					return "Yes"
				}
				return "No"
			}(k),
		})
	}

	return c, nil

}

func (c *ListController) Render() error {
	tableData := fctl.Map(c.store.Profiles, func(p *Profile) []string {
		return []string{
			p.Name,
			p.Active,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Name", "Active"})

	pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
	return nil
}

func NewListCommand() *cobra.Command {
	config := NewListConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
