package webhooks

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useListWebhook         = "list"
	descriptionListWebhook = "List all webhooks"
)

type ListWebhookStore struct {
	Webhooks []shared.WebhooksConfig `json:"webhooks"`
}

func NewDefaultListWebhookStore() *ListWebhookStore {
	return &ListWebhookStore{
		Webhooks: []shared.WebhooksConfig{},
	}
}

type ListWebhookControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewListWebhookControllerConfig() *ListWebhookControllerConfig {
	flags := flag.NewFlagSet(useListWebhook, flag.ExitOnError)
	fctl.WithGlobalFlags(flags)

	return &ListWebhookControllerConfig{
		context:     nil,
		use:         useListWebhook,
		description: descriptionListWebhook,
		aliases: []string{
			"ls", "l",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*ListWebhookStore] = (*ListWebhookController)(nil)

type ListWebhookController struct {
	store  *ListWebhookStore
	config ListWebhookControllerConfig
}

func NewListWebhookController(config ListWebhookControllerConfig) *ListWebhookController {
	return &ListWebhookController{
		store:  NewDefaultListWebhookStore(),
		config: config,
	}
}

func (c *ListWebhookController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *ListWebhookController) GetContext() context.Context {
	return c.config.context
}

func (c *ListWebhookController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *ListWebhookController) GetStore() *ListWebhookStore {
	return c.store
}

func (c *ListWebhookController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *ListWebhookController) Run() (fctl.Renderable, error) {

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

	webhookClient, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	request := operations.GetManyConfigsRequest{}
	response, err := webhookClient.Webhooks.GetManyConfigs(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "listing all config")
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Webhooks = response.ConfigsResponse.Cursor.Data

	return c, nil
}

func (c *ListWebhookController) Render() error {
	// TODO: WebhooksConfig is missing ?
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(c.config.out).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.Webhooks,
					func(src shared.WebhooksConfig) []string {
						return []string{
							src.ID,
							src.CreatedAt.Format(time.RFC3339),
							src.Secret,
							src.Endpoint,
							fctl.BoolToString(src.Active),
							strings.Join(src.EventTypes, ","),
						}
					}),
				[]string{"ID", "Created at", "Secret", "Endpoint", "Active", "Event types"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}
	return nil
}

func NewListCommand() *cobra.Command {

	config := NewListWebhookControllerConfig()

	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithController[*ListWebhookStore](NewListWebhookController(*config)),
	)
}
