package webhooks

import (
	"flag"
	"fmt"
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

func NewDeleteWebhookControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDeleteWebhook, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useDeleteWebhook,
		descriptionDeleteWebhook,
		[]string{
			"delete",
			"del",
		},
		os.Stdout,
		flags,
	)

}

var _ fctl.Controller[*DeleteWebhookStore] = (*DeleteWebhookController)(nil)

type DeleteWebhookController struct {
	store      *DeleteWebhookStore
	config     fctl.ControllerConfig
	fctlConfig *fctl.Config
}

func NewDeleteWebhookController(config fctl.ControllerConfig) *DeleteWebhookController {
	return &DeleteWebhookController{
		store:  NewDefaultDeleteWebhookStore(),
		config: config,
	}
}

func (c *DeleteWebhookController) GetStore() *DeleteWebhookStore {
	return c.store
}

func (c *DeleteWebhookController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *DeleteWebhookController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}
	c.fctlConfig = cfg

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

	if len(c.config.GetArgs()) == 0 {
		return nil, errors.New("missing config id")
	}

	request := operations.DeleteConfigRequest{
		ID: c.config.GetArgs()[0],
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
		pterm.Warning.WithShowLineNumber(false).Printfln("Config %s not found", c.config.GetArgs()[0])
		return nil
	}

	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.ErrorMessage)
		return nil
	}

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Config deleted successfully")

	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteWebhookControllerConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithAliases(config.GetAliases()...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteWebhookStore](NewDeleteWebhookController(*config)),
	)
}
