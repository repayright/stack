package webhooks

import (
	"flag"
	"fmt"
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

func NewListWebhookControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useListWebhook, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useListWebhook,
		descriptionListWebhook,
		descriptionListWebhook,
		[]string{
			"list",
			"ls",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*ListWebhookStore] = (*ListWebhookController)(nil)

type ListWebhookController struct {
	store  *ListWebhookStore
	config fctl.ControllerConfig
}

func NewListWebhookController(config fctl.ControllerConfig) *ListWebhookController {
	return &ListWebhookController{
		store:  NewDefaultListWebhookStore(),
		config: config,
	}
}

func (c *ListWebhookController) GetStore() *ListWebhookStore {
	return c.store
}

func (c *ListWebhookController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *ListWebhookController) Run() (fctl.Renderable, error) {

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
		WithWriter(c.config.GetOut()).
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

	return fctl.NewCommand(config.GetUse(),
		fctl.WithController[*ListWebhookStore](NewListWebhookController(*config)),
	)
}
