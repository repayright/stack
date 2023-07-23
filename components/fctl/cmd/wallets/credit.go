package wallets

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"
	"math/big"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	cobra "github.com/spf13/cobra"
)

const (
	useCredit   = "credit <amount> <asset>"
	shortCredit = "Credit a wallet"
	sourceFlag  = "source"
)

type CreditStore struct {
	Success bool `json:"success"`
}
type CreditController struct {
	store  *CreditStore
	config *config.ControllerConfig
}

var _ config.Controller[*CreditStore] = (*CreditController)(nil)

func NewCreditStore() *CreditStore {
	return &CreditStore{}
}

func NewCreditController(config *config.ControllerConfig) *CreditController {
	return &CreditController{
		store:  NewCreditStore(),
		config: config,
	}
}
func NewCreditConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useCredit, flag.ExitOnError)
	config.WithMetadataFlag(flags)
	config.WithConfirmFlag(flags)

	flags.String(balanceFlag, "", "Balance to credit")
	flags.String(sourceFlag, "", `Use --source account=<account> | --source wallet=id:<wallet-id>[/<balance>] | --source wallet=name:<wallet-name>[/<balance>]`)

	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)

	return config.NewControllerConfig(
		useCredit,
		shortCredit,
		shortCredit,
		[]string{
			"list",
			"ls",
		},
		flags,
		config.Organization, config.Stack,
	)
}
func (c *CreditController) GetStore() *CreditStore {
	return c.store
}

func (c *CreditController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *CreditController) Run() (config.Renderer, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to credit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) < 2 {
		return nil, errors.New("You need to specify amount and asset")
	}

	amountStr := c.config.GetArgs()[0]
	asset := c.config.GetArgs()[1]
	walletID, err := internal.RetrieveWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	if walletID == "" {
		return nil, errors.New("You need to specify wallet id using --id or --name flags")
	}

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	metadata, err := fctl.ParseMetadata(config.GetStringSlice(flags, config.M))
	if err != nil {
		return nil, err
	}

	sources := make([]shared.Subject, 0)
	for _, sourceStr := range config.GetStringSlice(flags, sourceFlag) {
		source, err := internal.ParseSubject(sourceStr, flags, ctx, client)
		if err != nil {
			return nil, err
		}
		sources = append(sources, *source)
	}

	request := operations.CreditWalletRequest{
		ID: walletID,
		CreditWalletRequest: &shared.CreditWalletRequest{
			Amount: shared.Monetary{
				Asset:  asset,
				Amount: amount,
			},
			Metadata: metadata,
			Sources:  sources,
			Balance:  formance.String(config.GetString(flags, balanceFlag)),
		},
	}
	response, err := client.Wallets.CreditWallet(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "crediting wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true

	return c, nil
}

func (c *CreditController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Wallet credited successfully!")
	return nil
}
func NewCreditWalletCommand() *cobra.Command {
	c := NewCreditConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController(NewCreditController(c)),
	)
}
