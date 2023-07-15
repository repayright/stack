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

func NewChangeSecretControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useChangeSecret, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	c := fctl.NewControllerConfig(
		useChangeSecret,
		descriptionChangeSecret,
		[]string{"cs"},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortDescriptionChangeSecret)

	return c
}

var _ fctl.Controller[*ChangeSecretStore] = (*ChangeSecretWebhookController)(nil)

type ChangeSecretWebhookController struct {
	store  *ChangeSecretStore
	config fctl.ControllerConfig
}

func NewChangeSecretWebhookController(config fctl.ControllerConfig) *ChangeSecretWebhookController {
	return &ChangeSecretWebhookController{
		store:  NewDefaultChangeSecretStore(),
		config: config,
	}
}

func (c *ChangeSecretWebhookController) GetStore() *ChangeSecretStore {
	return c.store
}

func (c *ChangeSecretWebhookController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *ChangeSecretWebhookController) Run() (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to change a webhook secret") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	secret := ""

	if len(c.config.GetArgs()) == 0 {
		return nil, fmt.Errorf("missing config-id")
	}

	if len(c.config.GetArgs()) > 1 {
		secret = c.config.GetArgs()[1]
	}

	response, err := client.Webhooks.
		ChangeConfigSecret(ctx, operations.ChangeConfigSecretRequest{
			ConfigChangeSecret: &shared.ConfigChangeSecret{
				Secret: secret,
			},
			ID: c.config.GetArgs()[0],
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
	pterm.Success.WithWriter(c.config.GetOut()).Printfln(
		"Config '%s' updated successfully with new secret", c.store.ID)
	return nil
}

func NewChangeSecretCommand() *cobra.Command {

	config := NewChangeSecretControllerConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(*config.GetShortDescription()),
		fctl.WithDescription(config.GetDescription()),
		fctl.WithAliases(config.GetAliases()...),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithController[*ChangeSecretStore](NewChangeSecretWebhookController(*config)),
	)
}
