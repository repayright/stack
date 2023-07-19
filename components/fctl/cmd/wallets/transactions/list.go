package transactions

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useListTransactions         = "list"
	descriptionListTransactions = "List all transactions"
)

type ListStore struct {
	Transactions []shared.WalletsTransaction `json:"transactions"`
}
type ListController struct {
	store  *ListStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useListTransactions, flag.ExitOnError)
	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)

	c := fctl.NewControllerConfig(
		useListTransactions,
		descriptionListTransactions,
		[]string{
			"ls",
			"l",
		},
		os.Stdout,
		flags,
	)

	return c
}
func NewListController(config fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewDefaultListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving organization ID in transaction list")
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, errors.Wrap(err, "resolving stack in transaction list")
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client in transaction list")
	}

	walletID, err := internal.RetrieveWalletID(flags, ctx, client)
	if err != nil {
		return nil, errors.Wrap(err, "retrieve wallet ID in transaction list")
	}

	request := operations.GetTransactionsRequest{
		WalletID: &walletID,
	}
	response, err := client.Wallets.GetTransactions(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "listing transactions")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transactions = response.GetTransactionsResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render() error {
	if len(c.store.Transactions) == 0 {
		fctl.Println("No transactions found.")
		return nil
	}

	tableData := fctl.Map(c.store.Transactions, func(tx shared.WalletsTransaction) []string {
		return []string{
			fmt.Sprintf("%d", tx.Txid),
			tx.Timestamp.Format(time.RFC3339),
			fctl.MetadataAsShortString(tx.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Date", "Metadata"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewListCommand() *cobra.Command {
	c := NewListConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithAliases(c.GetAliases()...),
		fctl.WithShortDescription(c.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(*c)),
	)
}
