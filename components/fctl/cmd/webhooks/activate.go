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
	useActivateWebhook         = "activate <config-id>"
	descriptionActivateWebhook = "Activate one config"
)

type ActivateWebhookStore struct {
	Success bool `json:"success"`
}

func NewDefaultVersionStore() *ActivateWebhookStore {
	return &ActivateWebhookStore{
		Success: true,
	}
}

type ActivateWebhookControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewActivateWebhookControllerConfig() *ActivateWebhookControllerConfig {
	flags := flag.NewFlagSet(useActivateWebhook, flag.ExitOnError)
	fctl.WithGlobalFlags(flags)
	fctl.WithConfirmFlag(flags)

	return &ActivateWebhookControllerConfig{
		context:     nil,
		use:         useActivateWebhook,
		description: descriptionActivateWebhook,
		aliases: []string{
			"ac", "a",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*ListWebhookStore] = (*ListWebhookController)(nil)

type ActivateWebhookController struct {
	store  *ActivateWebhookStore
	config ActivateWebhookControllerConfig
}

func NewActivateWebhookController(config ActivateWebhookControllerConfig) *ActivateWebhookController {
	return &ActivateWebhookController{
		store:  NewDefaultVersionStore(),
		config: config,
	}
}

func (c *ActivateWebhookController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *ActivateWebhookController) GetContext() context.Context {
	return c.config.context
}

func (c *ActivateWebhookController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *ActivateWebhookController) GetStore() *ActivateWebhookStore {
	return c.store
}

func (c *ActivateWebhookController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *ActivateWebhookController) Run() (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(flags, stack, "You are bout to activate a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.args) < 1 {
		return nil, fmt.Errorf("missing config id")
	}

	request := operations.ActivateConfigRequest{
		ID: c.config.args[0],
	}
	response, err := client.Webhooks.ActivateConfig(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "activating config")
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return c, nil
}

func (c *ActivateWebhookController) Render() error {
	pterm.Success.WithWriter(c.config.out).Printfln("Config activated successfully")

	return nil
}

func NewActivateCommand() *cobra.Command {
	config := NewActivateWebhookControllerConfig()
	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*ActivateWebhookStore](NewActivateWebhookController(*config)),
	)
}
