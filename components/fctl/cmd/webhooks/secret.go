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
	useChangeSecret              = "change-secret <config-id> <secret>"
	shortDescriptionChangeSecret = "Change the signing secret of a config"
	descriptionChangeSecret      = "Change the signing secret of a config. You can bring your own secret. If not passed or empty, a secret is automatically generated. The format is a string of bytes of size 24, base64 encoded. (larger size after encoding)"
)

type ChangeSecretStore struct {
	Secret string `json:"secret"`
	ID     string `json:"id"`
}

func NewDefaultChangeSecretStore() *ChangeSecretStore {
	return &ChangeSecretStore{
		Secret: "",
		ID:     "",
	}
}

type ChangeSecretControllerConfig struct {
	context          context.Context
	use              string
	shortDescription string
	description      string
	aliases          []string
	out              io.Writer
	flags            *flag.FlagSet
	args             []string
}

func NewChangeSecretControllerConfig() *ChangeSecretControllerConfig {
	flags := flag.NewFlagSet(useChangeSecret, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	fctl.WithGlobalFlags(flags)

	return &ChangeSecretControllerConfig{
		context:          nil,
		use:              useChangeSecret,
		shortDescription: shortDescriptionChangeSecret,
		description:      descriptionChangeSecret,
		aliases: []string{
			"cs",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*ChangeSecretStore] = (*ChangeSecretWebhookController)(nil)

type ChangeSecretWebhookController struct {
	store  *ChangeSecretStore
	config ChangeSecretControllerConfig
}

func NewChangeSecretWebhookController(config ChangeSecretControllerConfig) *ChangeSecretWebhookController {
	return &ChangeSecretWebhookController{
		store:  NewDefaultChangeSecretStore(),
		config: config,
	}
}

func (c *ChangeSecretWebhookController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *ChangeSecretWebhookController) GetContext() context.Context {
	return c.config.context
}

func (c *ChangeSecretWebhookController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *ChangeSecretWebhookController) GetStore() *ChangeSecretStore {
	return c.store
}

func (c *ChangeSecretWebhookController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *ChangeSecretWebhookController) Run() (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to change a webhook secret") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	secret := ""

	if len(c.config.args) == 0 {
		return nil, fmt.Errorf("missing config-id")
	}

	if len(c.config.args) > 1 {
		secret = c.config.args[1]
	}

	response, err := client.Webhooks.
		ChangeConfigSecret(ctx, operations.ChangeConfigSecretRequest{
			ConfigChangeSecret: &shared.ConfigChangeSecret{
				Secret: secret,
			},
			ID: c.config.args[0],
		})
	if err != nil {
		return nil, errors.Wrap(err, "changing secret")
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.ID = response.ConfigResponse.Data.ID
	c.store.Secret = response.ConfigResponse.Data.Secret

	return c, nil
}

func (c *ChangeSecretWebhookController) Render() error {
	pterm.Success.WithWriter(c.config.out).Printfln(
		"Config '%s' updated successfully with new secret", c.store.ID)
	return nil
}

func NewChangeSecretCommand() *cobra.Command {

	config := NewChangeSecretControllerConfig()

	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.shortDescription),
		fctl.WithDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*ChangeSecretStore](NewChangeSecretWebhookController(*config)),
	)
}
