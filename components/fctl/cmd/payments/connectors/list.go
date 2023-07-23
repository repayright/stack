package connectors

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	useList         = "list"
	descriptionList = "List all enabled connectors"
)

type ListStore struct {
	Connectors []shared.ConnectorsResponseData `json:"connectors"`
}

func NewListStore() *ListStore {
	return &ListStore{
		Connectors: []shared.ConnectorsResponseData{},
	}
}

func NewListConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)

	return config.NewControllerConfig(
		useList,
		descriptionList,
		"",
		[]string{
			"list",
			"ls",
		},
		flags,
		config.Organization, config.Stack,
	)

}

var _ config.Controller[*ListStore] = (*ListController)(nil)

func NewListController(config *config.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

type ListController struct {
	store  *ListStore
	config *config.ControllerConfig
}

func (c *ListController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run() (modelutils.Renderable, error) {

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

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	response, err := client.Payments.ListAllConnectors(ctx)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if response.ConnectorsResponse == nil {
		return nil, fmt.Errorf("unexpected response: %v", response)
	}

	c.store.Connectors = response.ConnectorsResponse.Data

	return c, nil
}

func (c *ListController) Render() error {
	tableData := fctl.Map(c.store.Connectors, func(connector shared.ConnectorsResponseData) []string {
		return []string{
			string(*connector.Provider),
			fctl.BoolToString(*connector.Enabled),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Provider", "Enabled"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {
	c := NewListConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithController[*ListStore](NewListController(c)),
	)
}
