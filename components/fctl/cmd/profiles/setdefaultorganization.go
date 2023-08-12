package profiles

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/cmd/profiles/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useSetOrg         = "set-default-organization <organization-id>"
	shortSetOrg       = "Set default organization"
	descriptionSetOrg = "Set default organization"
)

type SetOrgStore struct {
	Success bool `json:"success"`
}

func NewSetOrgStore() *SetOrgStore {
	return &SetOrgStore{
		Success: false,
	}
}

func NewSetOrgConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useSetOrg, flag.ExitOnError)

	c := config.NewControllerConfig(
		useSetOrg,
		descriptionSetOrg,
		shortSetOrg,
		[]string{
			"set-org",
			"sdo",
		},
		flags,
	)

	return c
}

var _ config.Controller = (*SetOrgController)(nil)

type SetOrgController struct {
	store  *SetOrgStore
	config *config.ControllerConfig
}

func NewSetOrgController(config *config.ControllerConfig) *SetOrgController {
	return &SetOrgController{
		store:  NewSetOrgStore(),
		config: config,
	}
}

func (c *SetOrgController) GetStore() any {
	return c.store
}
func (c *SetOrgController) GetKeyMapAction() *config.KeyMapHandler {
	return nil
}

func (c *SetOrgController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *SetOrgController) Run() (config.Renderer, error) {

	flags := c.config.GetAllFLags()
	args := flags.Args()

	if len(args) < 1 {
		return nil, errors.New("Please provide a profile name")
	}

	cfg, err := config.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	config.GetCurrentProfile(flags, cfg).SetDefaultOrganization(args[0])

	if err := cfg.Persist(); err != nil {
		return nil, errors.Wrap(err, "Updating config")
	}

	c.store.Success = true
	return c, nil
}

func (c *SetOrgController) Render() (tea.Model, error) {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Default organization updated!")
	return nil, nil
}

func NewSetDefaultOrganizationCommand() *cobra.Command {

	config := NewSetOrgConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController(NewSetOrgController(config)),
	)
}
