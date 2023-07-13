package webhooks

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDeleteWebhook         = "delete <config-id>"
	descriptionDeleteWebhook = "Delete a config"
)

type DeleteWebhookStore struct {
	ErrorResponse *shared.ErrorResponse `json:"error"`
	Success       bool                  `json:"success"`
}

func NewDefaultDeleteWebhookStore() *DeleteWebhookStore {
	return &DeleteWebhookStore{
		Success: true,
	}
}

type DeleteWebhookControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
	fctlConfig  *fctl.Config
}

func NewDeleteWebhookControllerConfig() *DeleteWebhookControllerConfig {
	flags := flag.NewFlagSet(useDeleteWebhook, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	fctl.WithGlobalFlags(flags)

	return &DeleteWebhookControllerConfig{
		context:     nil,
		use:         useDeleteWebhook,
		description: descriptionDeleteWebhook,
		aliases: []string{
			"rm", "del", "d",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*DeleteWebhookStore] = (*DeleteWebhookController)(nil)

type DeleteWebhookController struct {
	store  *DeleteWebhookStore
	config DeleteWebhookControllerConfig
}

func NewDeleteWebhookController(config DeleteWebhookControllerConfig) *DeleteWebhookController {
	return &DeleteWebhookController{
		store:  NewDefaultDeleteWebhookStore(),
		config: config,
	}
}

func (c *DeleteWebhookController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *DeleteWebhookController) GetContext() context.Context {
	return c.config.context
}

func (c *DeleteWebhookController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *DeleteWebhookController) GetStore() *DeleteWebhookStore {
	return c.store
}

func (c *DeleteWebhookController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *DeleteWebhookController) Run() (fctl.Renderable, error) {
	flags := c.config.flags
	ctx := c.config.context

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}
	c.config.fctlConfig = cfg

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to delete a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	webhookClient, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.args) == 0 {
		return nil, errors.New("missing config id")
	}

	request := operations.DeleteConfigRequest{
		ID: c.config.args[0],
	}
	response, err := webhookClient.Webhooks.DeleteConfig(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "deleting config")
	}

	if response.ErrorResponse != nil {
		if response.ErrorResponse.ErrorCode == "NOT_FOUND" {
			c.store.ErrorResponse = response.ErrorResponse
			c.store.Success = false
			return c, nil
		}
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = response.StatusCode == 200

	return c, nil
}

func (c *DeleteWebhookController) Render() error {
	if !c.store.Success {
		pterm.Warning.WithShowLineNumber(false).Printfln("Config %s not found", c.config.args[0])
		return nil
	}

	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.ErrorMessage)
		return nil
	}

	pterm.Success.WithWriter(c.config.out).Printfln("Config deleted successfully")

	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteWebhookControllerConfig()
	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*DeleteWebhookStore](NewDeleteWebhookController(*config)),
	)
}
