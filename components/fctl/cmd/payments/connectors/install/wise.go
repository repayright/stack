package install

import (
	"flag"
	"fmt"
	"os"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useWiseConnector              = internal.WiseConnector + " <api-key>"
	descriptionWiseConnector      = "Install Wise connector"
	shortDescriptionWiseConnector = "Install Wise connector"
)

type WiseStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewWiseStore() *WiseStore {
	return &WiseStore{
		Success: false,
	}
}
func NewWiseConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useWiseConnector, flag.ExitOnError)
	flags.String(PollingPeriodFlag, DefaultPollingPeriod, "Polling duration")
	c := fctl.NewControllerConfig(
		useWiseConnector,
		descriptionWiseConnector,
		[]string{},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortDescriptionWiseConnector)

	return c
}

var _ fctl.Controller[*WiseStore] = (*WiseController)(nil)

type WiseController struct {
	store  *WiseStore
	config fctl.ControllerConfig
}

func NewWiseController(config fctl.ControllerConfig) *WiseController {
	return &WiseController{
		store:  NewWiseStore(),
		config: config,
	}
}

func (c *WiseController) GetStore() *WiseStore {
	return c.store
}

func (c *WiseController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *WiseController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	if len(args) < 1 {
		return nil, fmt.Errorf("missing api key")
	}

	soc, err := fctl.GetStackOrganizationConfigApprobation(flags, ctx, "You are about to install connector '%s'", internal.WiseConnector)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	response, err := paymentsClient.Payments.InstallConnector(ctx, operations.InstallConnectorRequest{
		RequestBody: shared.WiseConfig{
			APIKey:        args[1],
			PollingPeriod: fctl.Ptr(fctl.GetString(flags, PollingPeriodFlag)),
		},
		Connector: shared.ConnectorWise,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.WiseConnector

	return c, nil
}

func (c *WiseController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
func NewWiseCommand() *cobra.Command {
	config := NewWiseConfig()
	c := NewWiseController(*config)
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*WiseStore](c),
	)
}
