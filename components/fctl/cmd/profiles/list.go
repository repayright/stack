package profiles

import (
	"flag"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/ui"

	fctl "github.com/formancehq/fctl/pkg"
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

var _ config.Controller = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *config.ControllerConfig
	keymap *config.KeyMapHandler
}

func NewListController(conf *config.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: conf,
		keymap: config.NewKeyMapHandler(),
	}
}
func (c *ListController) GetKeyMapAction() *config.KeyMapHandler {
	return c.keymap
}
func (c *ListController) GetStore() any {
	return c.store
}

func (c *ListController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (config.Renderer, error) {

	flags := c.config.GetAllFLags()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profiles := cfg.GetProfiles()
	p := fctl.MapKeys(profiles)

	currentProfileName := fctl.GetCurrentProfileName(flags, cfg)

	c.store.Profiles = make([]*Profile, len(p))
	for i, k := range p {
		c.store.Profiles[i] = &Profile{
			Name: k,
			Active: func(k string) string {
				if currentProfileName == k {
					return "Yes"
				}
				return "No"
			}(k),
		}
	}

	return c, nil

}

func (c *ListController) Render() (tea.Model, error) {
	flags := c.config.GetAllFLags()
	tableData := fctl.Map(c.store.Profiles, func(p *Profile) table.Row {
		return []string{
			p.Name,
			p.Active,
		}
	})

	columns := ui.NewArrayColumn(
		ui.NewColumn("Name", 10),
		ui.NewColumn("Active", 20),
	)

	opts := ui.NewTableOptions(columns, tableData)

	if config.GetString(flags, config.OutputFlag) == "plain" {
		opt := ui.WithHeight(len(tableData))
		// Add Deleted At column if --deleted flag is set
		return ui.NewTableModel(columns, append(opts, opt)...), nil
	}

	opts = ui.NewTableOptions(ui.WithFullScreenTable(columns), tableData)

	return ui.NewTableModel(columns, opts...), nil
}

func NewListCommand() *cobra.Command {
	config := NewListConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithController(NewListController(config)),
	)
}
