package ledger

import (
	"flag"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

const (
	useList         = "list"
	shortList       = "List ledgers"
	descriptionList = "List all ledgers"
)

type ListStore struct {
	Ledgers []string `json:"ledgers"`
}

func NewListStore() *ListStore {
	return &ListStore{
		Ledgers: []string{},
	}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useList,
		descriptionList,
		shortList,
		[]string{
			"l", "ls",
		},
		os.Stdout,
		flags,
	)
}

type ListController struct {
	store  *ListStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

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

	response, err := ledgerClient.Ledger.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	c.store.Ledgers = response.ConfigInfoResponse.Data.Config.Storage.Ledgers

	return c, nil
}

func (c *ListController) Render() error {
	tableData := fctl.Map(c.store.Ledgers, func(ledger string) []string {
		return []string{
			ledger,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Name"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}
func NewListCommand() *cobra.Command {
	config := NewListConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(*config)),
	)
}
