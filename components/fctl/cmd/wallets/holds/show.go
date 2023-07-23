package holds

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	"github.com/formancehq/fctl/cmd/wallets/internal/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	useShow   = "show <hold-id>"
	shortShow = "Show a hold"
)

type ShowStore struct {
	Hold shared.ExpandedDebitHold `json:"hold"`
}
type ShowController struct {
	store  *ShowStore
	config *config.ControllerConfig
}

var _ config.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowStore() *ShowStore {
	return &ShowStore{}
}
func NewShowConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)

	return config.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh",
		},
		flags,
		config.Organization, config.Stack,
	)

}
func NewShowController(config *config.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", len(c.config.GetArgs()))
	}

	request := operations.GetHoldRequest{
		HoldID: c.config.GetArgs()[0],
	}
	response, err := client.Wallets.GetHold(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "getting hold")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Hold = response.GetHoldResponse.Data

	return c, nil
}

func (c *ShowController) Render() error {

	return views.PrintHold(c.config.GetOut(), c.store.Hold)

}
func NewShowCommand() *cobra.Command {

	c := NewShowConfig()

	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(c)),
	)
}
