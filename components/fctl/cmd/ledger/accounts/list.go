package accounts

import (
	"flag"
	"fmt"
	fctlConfig "github.com/formancehq/fctl/pkg/config"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList         = "list"
	shortList       = "List accounts"
	descriptionList = "List all accounts"
)

type ListStore struct {
	Accounts []shared.Account `json:"accounts"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.String(fctlConfig.MetadataFlag, "", "Filter accounts with metadata")

	return fctlConfig.NewControllerConfig(
		useList,
		descriptionList,
		shortList,
		[]string{
			"l", "ls",
		},
		flags, fctlConfig.Organization, fctlConfig.Stack, fctlConfig.Ledger,
	)
}

type ListController struct {
	store  *ListStore
	config *fctlConfig.ControllerConfig
}

var _ fctlConfig.Controller[*ListStore] = (*ListController)(nil)

func NewListController(config *fctlConfig.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctlConfig.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctlConfig.Renderer, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	metadata, err := fctl.ParseMetadata(fctlConfig.GetStringSlice(flags, fctlConfig.MetadataFlag))
	if err != nil {
		return nil, err
	}

	request := operations.ListAccountsRequest{
		Ledger:   fctlConfig.GetString(flags, internal.LedgerFlag),
		Metadata: metadata,
	}
	rsp, err := ledgerClient.Ledger.ListAccounts(ctx, request)
	if err != nil {
		return nil, err
	}

	if rsp.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", rsp.ErrorResponse.ErrorCode, rsp.ErrorResponse.ErrorMessage)
	}

	if rsp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	c.store.Accounts = rsp.AccountsCursorResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render() error {

	tableData := fctl.Map(c.store.Accounts, func(account shared.Account) []string {
		return []string{
			account.Address,
			fctl.MetadataAsShortString(account.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Address", "Metadata"})
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
		fctl.WithController[*ListStore](NewListController(c)),
	)
}
