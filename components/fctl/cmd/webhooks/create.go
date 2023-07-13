package webhooks

import (
	"context"
	"flag"
	"fmt"
	"io"
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

type CreateWebhookStore struct {
	Webhook shared.WebhooksConfig `json:"webhook"`
}

func NewDefaultCreateWebhookStore() *CreateWebhookStore {
	return &CreateWebhookStore{
		Webhook: shared.WebhooksConfig{},
	}
}

type CreateWebhookControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewCreateWebhookControllerConfig() *CreateWebhookControllerConfig {
	flags := flag.NewFlagSet(useCreateWebhook, flag.ExitOnError)
	flags.String(secretFlag, "", "Bring your own webhooks signing secret. If not passed or empty, a secret is automatically generated. The format is a string of bytes of size 24, base64 encoded. (larger size after encoding)")
	fctl.WithGlobalFlags(flags)
	fctl.WithConfirmFlag(flags)

	return &CreateWebhookControllerConfig{
		context:     nil,
		use:         useCreateWebhook,
		description: descriptionCreateWebhook,
		aliases: []string{
			"cr",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*CreateWebhookStore] = (*CreateWebhookController)(nil)

type CreateWebhookController struct {
	store  *CreateWebhookStore
	config CreateWebhookControllerConfig
}

func NewCreateWebhookController(config CreateWebhookControllerConfig) *CreateWebhookController {
	return &CreateWebhookController{
		store:  NewDefaultCreateWebhookStore(),
		config: config,
	}
}

func (c *CreateWebhookController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *CreateWebhookController) GetContext() context.Context {
	return c.config.context
}

func (c *CreateWebhookController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *CreateWebhookController) GetStore() *CreateWebhookStore {
	return c.store
}

func (c *CreateWebhookController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *CreateWebhookController) Run() (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.args) < 2 {
		return nil, fmt.Errorf("at least one event type is required")
	}

	if _, err := url.Parse(c.config.args[0]); err != nil {
		return nil, errors.Wrap(err, "invalid endpoint URL")
	}

	secret := fctl.GetString(flags, secretFlag)

	response, err := client.Webhooks.InsertConfig(ctx, shared.ConfigUser{
		Endpoint:   c.config.args[0],
		EventTypes: c.config.args[1:],
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
	pterm.Success.WithWriter(c.config.out).Printfln("Config created successfully")
	return nil
}

func NewCreateCommand() *cobra.Command {
	config := NewCreateWebhookControllerConfig()
	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*CreateWebhookStore](NewCreateWebhookController(*config)),
	)
}
