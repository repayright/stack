package webhooks

import (
	"context"
	"flag"
	"fmt"
	"io"
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

type DesactivateWebhookControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewDesactivateWebhookControllerConfig() *DesactivateWebhookControllerConfig {
	flags := flag.NewFlagSet(useDesactivateWebhook, flag.ExitOnError)
	fctl.WithGlobalFlags(flags)
	fctl.WithConfirmFlag(flags)

	return &DesactivateWebhookControllerConfig{
		context:     nil,
		use:         useDesactivateWebhook,
		description: descriptionDesactivateWebhook,
		aliases: []string{
			"deac",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*DesactivateWebhookStore] = (*DesactivateWebhookController)(nil)

type DesactivateWebhookController struct {
	store  *DesactivateWebhookStore
	config DesactivateWebhookControllerConfig
}

func NewDesactivateWebhookController(config DesactivateWebhookControllerConfig) *DesactivateWebhookController {
	return &DesactivateWebhookController{
		store:  NewDefaultDesactivateWebhookStore(),
		config: config,
	}
}

func (c *DesactivateWebhookController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *DesactivateWebhookController) GetContext() context.Context {
	return c.config.context
}

func (c *DesactivateWebhookController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *DesactivateWebhookController) GetStore() *DesactivateWebhookStore {
	return c.store
}

func (c *DesactivateWebhookController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *DesactivateWebhookController) Run() (fctl.Renderable, error) {

	flags := c.config.flags
	ctx := c.config.context

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

	if len(c.config.args) < 1 {
		return nil, fmt.Errorf("missing config id")
	}

	request := operations.DeactivateConfigRequest{
		ID: c.config.args[0],
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

	pterm.Success.WithWriter(c.config.out).Printfln("Config deactivated successfully")

	return nil
}

func NewDeactivateCommand() *cobra.Command {

	config := NewDesactivateWebhookControllerConfig()

	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*DesactivateWebhookStore](NewDesactivateWebhookController(*config)),
	)
}
