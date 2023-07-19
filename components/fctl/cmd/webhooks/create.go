package webhooks

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	secretFlag               = "secret"
	useCreateWebhook         = "create <endpoint> [<event-type>...]"
	descriptionCreateWebhook = "Create a new config. At least one event type is required."
)

type CreateStore struct {
	Webhook shared.WebhooksConfig `json:"webhook"`
}

func NewDefaultCreateWebhookStore() *CreateStore {
	return &CreateStore{
		Webhook: shared.WebhooksConfig{},
	}
}

func NewCreateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreateWebhook, flag.ExitOnError)
	flags.String(secretFlag, "", "Bring your own webhooks signing secret. If not passed or empty, a secret is automatically generated. The format is a string of bytes of size 24, base64 encoded. (larger size after encoding)")
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useCreateWebhook,
		descriptionCreateWebhook,
		[]string{"cr"},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*CreateStore] = (*CreateWebhookController)(nil)

type CreateWebhookController struct {
	store  *CreateStore
	config fctl.ControllerConfig
}

func NewCreateWebhookController(config fctl.ControllerConfig) *CreateWebhookController {
	return &CreateWebhookController{
		store:  NewDefaultCreateWebhookStore(),
		config: config,
	}
}

func (c *CreateWebhookController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateWebhookController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *CreateWebhookController) Run() (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) < 2 {
		return nil, fmt.Errorf("at least one event type is required")
	}

	if _, err := url.Parse(c.config.GetArgs()[0]); err != nil {
		return nil, errors.Wrap(err, "invalid endpoint URL")
	}

	secret := fctl.GetString(flags, secretFlag)

	response, err := client.Webhooks.InsertConfig(ctx, shared.ConfigUser{
		Endpoint:   c.config.GetArgs()[0],
		EventTypes: c.config.GetArgs()[1:],
		Secret:     &secret,
	})

	if err != nil {
		return nil, errors.Wrap(err, "creating config")
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Webhook = response.ConfigResponse.Data

	return c, nil
}

func (c *CreateWebhookController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Config created successfully")
	return nil
}

func NewCreateCommand() *cobra.Command {
	config := NewCreateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithAliases(config.GetAliases()...),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithGoFlagSet(config.GetFlags()),
		fctl.WithController[*CreateStore](NewCreateWebhookController(*config)),
	)
}
