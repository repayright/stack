package balances

import (
	"flag"
	"fmt"
	"os"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	"github.com/formancehq/fctl/cmd/wallets/internal/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	useShow         = "show <balance-name>"
	descriptionShow = "Show a balance"
)

type ShowStore struct {
	Balance shared.BalanceWithAssets `json:"balance"`
}
type ShowController struct {
	store  *ShowStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{}
}
func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	internal.WithTargetingWalletByID(flags)
	internal.WithTargetingWalletByName(flags)
	c := fctl.NewControllerConfig(
		useShow,
		descriptionShow,
		[]string{
			"sh",
		},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(descriptionShow)

	return c
}
func NewShowController(config fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewDefaultShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	walletID, err := internal.RequireWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	if len(c.config.GetArgs()) != 1 {
		return nil, fmt.Errorf("invalid number of arguments")
	}

	request := operations.GetBalanceRequest{
		ID:          walletID,
		BalanceName: c.config.GetArgs()[0],
	}
	response, err := client.Wallets.GetBalance(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "getting balance")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Balance = response.GetBalanceResponse.Data

	return c, nil
}

func (c *ShowController) Render() error {
	return views.PrintBalance(c.config.GetOut(), c.store.Balance)
}

func NewShowCommand() *cobra.Command {
	c := NewShowConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithShortDescription(c.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(*c)),
	)
}
