package webhooks

import (
	"flag"
	"fmt"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDesactivateWebhook         = "deactivate <config-id>"
	descriptionDesactivateWebhook = "Deactivate one config"
)

type DesactivateWebhookStore struct {
	Success bool `json:"success"`
}

func NewDefaultDesactivateWebhookStore() *DesactivateWebhookStore {
	return &DesactivateWebhookStore{
		Success: true,
	}
}

func NewDesactivateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDesactivateWebhook, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useDesactivateWebhook,
		descriptionDesactivateWebhook,
		[]string{
			"deactivate",
			"des",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*DesactivateWebhookStore] = (*DesactivateWebhookController)(nil)

type DesactivateWebhookController struct {
	store  *DesactivateWebhookStore
	config fctl.ControllerConfig
}

func NewDesactivateWebhookController(config fctl.ControllerConfig) *DesactivateWebhookController {
	return &DesactivateWebhookController{
		store:  NewDefaultDesactivateWebhookStore(),
		config: config,
	}
}

func (c *DesactivateWebhookController) GetStore() *DesactivateWebhookStore {
	return c.store
}

func (c *DesactivateWebhookController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *DesactivateWebhookController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to deactivate a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) < 1 {
		return nil, fmt.Errorf("missing config id")
	}

	request := operations.DeactivateConfigRequest{
		ID: c.config.GetArgs()[0],
	}
	response, err := client.Webhooks.DeactivateConfig(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "deactivating config")
	}

	c.store.Success = !response.ConfigResponse.Data.Active

	// Check if there is an error
	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	// Check if the status code is >= 300
	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return c, nil
}

func (c *DesactivateWebhookController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Config deactivated successfully")

	return nil
}

func NewDeactivateCommand() *cobra.Command {

	config := NewDesactivateConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithAliases(config.GetAliases()...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DesactivateWebhookStore](NewDesactivateWebhookController(*config)),
	)
}
