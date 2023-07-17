package connectors

import (
	"flag"
	"fmt"
	"os"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	PaymentsConnectorsUninstall = "develop"
	useUninstall                = "uninstall <connector-name>"
	descriptionUninstall        = "Uninstall a connector"
	shortDescriptionUninstall   = "Uninstall a connector"
)

type UninstallStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewUninstallConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useUninstall, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	c := fctl.NewControllerConfig(
		useUninstall,
		descriptionUninstall,
		[]string{
			"uninstall", "u", "un",
		},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortDescriptionUninstall)

	return c
}

func NewUninstallStore() *UninstallStore {
	return &UninstallStore{
		Success:       false,
		ConnectorName: "",
	}
}

var _ fctl.Controller[*UninstallStore] = (*UninstallController)(nil)

type UninstallController struct {
	store  *UninstallStore
	config fctl.ControllerConfig
}

func NewUninstallController(config fctl.ControllerConfig) *UninstallController {
	return &UninstallController{
		store:  NewUninstallStore(),
		config: config,
	}
}

func (c *UninstallController) GetStore() *UninstallStore {
	return c.store
}

func (c *UninstallController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *UninstallController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to uninstall connector '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := client.Payments.UninstallConnector(ctx, operations.UninstallConnectorRequest{
		Connector: shared.Connector(args[0]),
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = args[0]

	return c, nil
}

func (c *UninstallController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector '%s' uninstalled!", c.store.ConnectorName)
	return nil
}

func NewUninstallCommand() *cobra.Command {

	c := NewUninstallConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgs(internal.AllConnectors...),
		fctl.WithShortDescription(*c.GetShortDescription()),
		fctl.WithController[*UninstallStore](NewUninstallController(*c)),
	)
}
