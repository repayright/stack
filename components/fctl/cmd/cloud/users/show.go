package users

import (
	"flag"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useShow   = "show <user-id>"
	shortShow = "Show user by id"
)

type ShowStore struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	return config.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh", "s",
		},
		flags,
		config.Organization,
	)
}

var _ config.Controller[*ShowStore] = (*ShowController)(nil)

type ShowController struct {
	store  *ShowStore
	config *config.ControllerConfig
}

func NewShowController(config *config.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *config.ControllerConfig {
	return c.config
}
func (c *ShowController) Run() (modelutils.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	userResponse, _, err := apiClient.DefaultApi.ReadUser(ctx, organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Id = userResponse.Data.Id
	c.store.Email = userResponse.Data.Email

	return c, nil
}

func (c *ShowController) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Id})
	tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.Email})

	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewShowCommand() *cobra.Command {
	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
