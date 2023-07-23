package instances

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useStop   = "stop <instance-id>"
	shortStop = "Stop a specific workflow instance"
)

type StopStore struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instanceId"`
}

func NewStopConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useStop, flag.ExitOnError)

	c := config.NewControllerConfig(
		useStop,
		shortStop,
		shortStop,
		[]string{},
		flags,
		config.Organization, config.Stack,
	)

	return c
}

type StopController struct {
	store  *StopStore
	config *config.ControllerConfig
}

var _ config.Controller[*StopStore] = (*StopController)(nil)

func NewStopStore() *StopStore {
	return &StopStore{}
}

func NewStopController(config *config.ControllerConfig) *StopController {
	return &StopController{
		store:  NewStopStore(),
		config: config,
	}
}

func (c *StopController) GetStore() *StopStore {
	return c.store
}

func (c *StopController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *StopController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()
	args := c.config.GetArgs()
	ctx := c.config.GetContext()

	soc, err := fctl.GetStackOrganizationConfig(flags, ctx, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.CancelEvent(ctx, operations.CancelEventRequest{
		InstanceID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.InstanceID = args[0]

	return c, nil
}

func (c *StopController) Render() error {
	// Print the instance information
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Workflow Instance with ID: %s successfully canceled ", c.store.InstanceID)

	return nil
}

func NewStopCommand() *cobra.Command {
	config := NewStopConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*StopStore](NewStopController(config)),
	)
}
