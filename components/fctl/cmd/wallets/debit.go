package wallets

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDebitWallet         = "debit <amount> <asset>"
	descriptionDebitWallet = "Debit a wallet"
	pendingFlag            = "pending"
	metadataFlag           = "metadata"
	descriptionFlag        = "description"
	balanceFlag            = "balance"
	destinationFlag        = "destination"
)

type DebitWalletStore struct {
	HoldID  *string `json:"holdId"`
	Success bool    `json:"success"`
}
type DebitWalletController struct {
	store  *DebitWalletStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*DebitWalletStore] = (*DebitWalletController)(nil)

func NewDefaultDebitWalletStore() *DebitWalletStore {
	return &DebitWalletStore{
		HoldID:  nil,
		Success: false,
	}
}
func NewDebitConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDebitWallet, flag.ExitOnError)
	flags.String(descriptionFlag, "", "Debit description")
	flags.String(pendingFlag, "", "Create a pending debit")
	flags.String(balanceFlag, "", "Balance to debit")
	flags.String(destinationFlag, "",
		`Use --destination account=<account> | --destination wallet=id:<wallet-id>[/<balance>] | --destination wallet=name:<wallet-name>[/<balance>]`)
	fctl.WithMetadataFlag(flags)
	fctl.WithConfirmFlag(flags)
	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)

	c := fctl.NewControllerConfig(
		useDebitWallet,
		descriptionDebitWallet,
		[]string{
			"deb",
		},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(descriptionDebitWallet)

	return c
}
func NewDebitWalletController(config fctl.ControllerConfig) *DebitWalletController {
	return &DebitWalletController{
		store:  NewDefaultDebitWalletStore(),
		config: config,
	}
}

func (c *DebitWalletController) GetStore() *DebitWalletStore {
	return c.store
}

func (c *DebitWalletController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *DebitWalletController) Run() (fctl.Renderable, error) {
	flags := c.config.GetFlags()
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to debit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	pending := fctl.GetBool(flags, pendingFlag)

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, metadataFlag))
	if err != nil {
		return nil, err
	}

	if len(c.config.GetArgs()) < 2 {
		return nil, errors.New("missing amount and asset")
	}

	amountStr := c.config.GetArgs()[0]
	asset := c.config.GetArgs()[1]
	walletID, err := internal.RequireWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	description := fctl.GetString(flags, descriptionFlag)

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	var destination *shared.Subject
	if destinationStr := fctl.GetString(flags, destinationFlag); destinationStr != "" {
		destination, err = internal.ParseSubject(destinationStr, flags, ctx, client)
		if err != nil {
			return nil, err
		}
	}

	response, err := client.Wallets.DebitWallet(ctx, operations.DebitWalletRequest{
		DebitWalletRequest: &shared.DebitWalletRequest{
			Amount: shared.Monetary{
				Asset:  asset,
				Amount: amount,
			},
			Pending:     &pending,
			Metadata:    metadata,
			Description: &description,
			Destination: destination,
			Balances:    fctl.GetStringSlice(flags, balanceFlag),
		},
		ID: walletID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "debiting wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if response.DebitWalletResponse != nil {
		c.store.HoldID = &response.DebitWalletResponse.Data.ID
	}

	c.store.Success = true

	return c, nil
}

func (c *DebitWalletController) Render() error {
	if c.store.HoldID != nil && *c.store.HoldID != "" {
		pterm.Success.WithWriter(c.config.GetOut()).Printfln("Wallet debited successfully with hold id '%s'!", *c.store.HoldID)
	} else {
		pterm.Success.WithWriter(c.config.GetOut()).Printfln("Wallet debited successfully!")
	}

	return nil

}

func NewDebitWalletCommand() *cobra.Command {
	c := NewDebitConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithShortDescription(*c.GetShortDescription()),
		fctl.WithAliases(c.GetAliases()...),
		fctl.WithArgs(cobra.RangeArgs(2, 3)),
		fctl.WithGoFlagSet(c.GetFlags()),
		fctl.WithController[*DebitWalletStore](NewDebitWalletController(*c)),
	)
}
