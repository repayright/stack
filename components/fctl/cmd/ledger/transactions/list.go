package transactions

import (
	"flag"
	"fmt"
	"os"
	"time"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	pageSizeFlag    = "page-size"
	accountFlag     = "account"
	destinationFlag = "dst"
	sourceFlag      = "src"
	endTimeFlag     = "end"
	startTimeFlag   = "start"
)

const (
	useList   = "list"
	shortList = "List transactions"
)

type ListStore struct {
	Transaction shared.TransactionsCursorResponseCursor `json:"transactionCursor"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.String(accountFlag, "", "Filter on account")
	flags.String(destinationFlag, "", "Filter on destination account")
	flags.String(endTimeFlag, "", "Consider transactions before date")
	flags.String(startTimeFlag, "", "Consider transactions after date")
	flags.String(sourceFlag, "", "Filter on source account")
	flags.String(internal.ReferenceFlag, "", "Filter on reference")
	flags.String(internal.MetadataFlag, "", "Filter transactions with metadata") //fctl.WithHiddenFlag(metadataFlag)
	flags.Int(pageSizeFlag, 5, "Page size")

	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"l", "ls",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config fctl.ControllerConfig
}

func NewListController(config fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
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
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, err
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, internal.MetadataFlag))
	if err != nil {
		return nil, err
	}

	var (
		endTime   time.Time
		startTime time.Time
	)
	if startTimeStr := fctl.GetString(flags, startTimeFlag); startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339Nano, startTimeStr)
		if err != nil {
			return nil, err
		}
	}
	if endTimeStr := fctl.GetString(flags, endTimeFlag); endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339Nano, endTimeStr)
		if err != nil {
			return nil, err
		}
	}

	ledger := fctl.GetString(flags, internal.LedgerFlag)
	response, err := ledgerClient.Ledger.ListTransactions(
		ctx,
		operations.ListTransactionsRequest{
			Account:     fctl.Ptr(fctl.GetString(flags, accountFlag)),
			Destination: fctl.Ptr(fctl.GetString(flags, destinationFlag)),
			EndTime:     &endTime,
			Ledger:      ledger,
			Metadata:    metadata,
			PageSize:    fctl.Ptr(int64(fctl.GetInt(flags, pageSizeFlag))),
			Reference:   fctl.Ptr(fctl.GetString(flags, internal.ReferenceFlag)),
			Source:      fctl.Ptr(fctl.GetString(flags, sourceFlag)),
			StartTime:   &startTime,
		},
	)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transaction = response.TransactionsCursorResponse.Cursor

	return c, nil
}

func (c *ListController) Render() error {
	if len(c.store.Transaction.Data) == 0 {
		fctl.Println("No transactions found.")
		return nil
	}

	tableData := fctl.Map(c.store.Transaction.Data, func(tx shared.ExpandedTransaction) []string {
		return []string{
			fmt.Sprintf("%d", tx.Txid),
			func() string {
				if tx.Reference == nil {
					return ""
				}
				return *tx.Reference
			}(),
			tx.Timestamp.Format(time.RFC3339),
			fctl.MetadataAsShortString(tx.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Reference", "Date", "Metadata"})

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {
	c := NewListConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(*c)),
	)
}
