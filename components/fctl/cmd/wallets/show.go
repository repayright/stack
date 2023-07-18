package wallets

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
	useShow   = "show"
	shortShow = "Show a wallet"
)

type ShowStore struct {
	Wallet shared.WalletWithBalances `json:"wallet"`
}
type ShowController struct {
	store  *ShowStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{
		Wallet: shared.WalletWithBalances{},
	}
}

func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	internal.WithTargetingWalletByID(flags)
	internal.WithTargetingWalletByName(flags)
	return fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh",
		},
		os.Stdout,
		flags,
	)
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

	walletID, err := internal.RetrieveWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}
	if walletID == "" {
		return nil, errors.New("You need to specify wallet id using --id or --name flags")
	}

	response, err := client.Wallets.GetWallet(ctx, operations.GetWalletRequest{
		ID: walletID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Wallet = response.GetWalletResponse.Data

	return c, nil
}

func (c *ShowController) Render() error {
	return views.PrintWalletWithMetadata(c.config.GetOut(), c.store.Wallet)
}

func NewShowCommand() *cobra.Command {
	c := NewShowConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithShortDescription(c.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ShowStore](NewShowController(*c)),
	)
}
