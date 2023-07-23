package stack

import (
	"flag"
	"github.com/charmbracelet/bubbles/table"
	"github.com/formancehq/fctl/pkg/ui"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// This defines the minimum length of the columns in the table
const (
	minLengthOrganizationId = 15
	minLengthStackId        = 8
	minLengthStackName      = 10
	minLengthApiUrl         = 48
	minLengthStackRegion    = 30
	minLengthStackCreatedAt = 20
	minLengthStackDeletedAt = 20
)
const (
	deletedFlag = "deleted"
	useList     = "list"
	shortList   = "List stacks"
)

type Stack struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Dashboard string  `json:"dashboard"`
	RegionID  string  `json:"region"`
	DeletedAt *string `json:"deletedAt"`
	CreatedAt *string `json:"createdAt"`
}

type ListStore struct {
	Stacks []Stack `json:"stacks"`
}

func NewListStore() *ListStore {
	return &ListStore{
		Stacks: []Stack{},
	}
}

func NewListControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.Bool(deletedFlag, false, "Show deleted stacks")

	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"list",
			"ls",
		},
		flags,
		fctl.Organization,
	)
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store        *ListStore
	profile      *fctl.Profile
	config       *fctl.ControllerConfig
	organization string
}

func NewListController(config *fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).
		Deleted(fctl.GetBool(flags, deletedFlag)).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "listing stacks")
	}

	c.profile = profile
	if len(rsp.Data) == 0 {
		return c, nil
	}

	c.organization = organization
	c.store.Stacks = fctl.Map(rsp.Data, func(stack membershipclient.Stack) Stack {
		return Stack{
			Id:        stack.Id,
			Name:      stack.Name,
			Dashboard: c.profile.ServicesBaseUrl(&stack).String(),
			RegionID:  stack.RegionID,
			CreatedAt: func() *string {
				if stack.CreatedAt != nil {
					t := stack.CreatedAt.Format(time.RFC3339)
					return &t
				}
				return nil
			}(),
			DeletedAt: func() *string {
				if stack.DeletedAt != nil {
					t := stack.DeletedAt.Format(time.RFC3339)
					return &t
				}
				return nil
			}(),
		}
	})

	return c, nil
}

func (c *ListController) Render() (ui.Model, error) {

	flags := c.config.GetAllFLags()

	// Create table rows
	tableData := fctl.Map(c.store.Stacks, func(stack Stack) table.Row {
		data := []string{
			c.organization,
			stack.Id,
			stack.Name,
			stack.Dashboard,
			stack.RegionID,
			*stack.CreatedAt,
		}

		if fctl.GetBool(flags, deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, *stack.DeletedAt)
			} else {
				data = append(data, "")
			}
		}

		return data
	})

	var columns ui.ArrayColumn

	// Add plain table option if --plain flag is set
	isPlain := fctl.GetString(flags, fctl.OutputFlag) == "plain"
	// Default Columns
	columns = ui.NewArrayColumn(
		ui.NewColumn("Organization Id", minLengthOrganizationId),
		ui.NewColumn("Stack Id", minLengthStackId),
		ui.NewColumn("Name", minLengthStackName),
		ui.NewColumn("API URL", minLengthApiUrl),
		ui.NewColumn("Region", minLengthStackRegion),
		ui.NewColumn("Created At", minLengthStackCreatedAt),
	)
	if fctl.GetBool(flags, deletedFlag) {
		columns = columns.AddColumn("Deleted At", minLengthStackDeletedAt)
	}
	// Default table options
	opts := ui.NewTableOptions(columns, tableData)
	if isPlain {
		opt := ui.WithHeight(len(tableData))
		// Add Deleted At column if --deleted flag is set
		return ui.NewTableModel(columns, append(opts, opt)...), nil
	}

	opts = ui.NewTableOptions(ui.WithFullScreenTable(columns), tableData)

	return ui.NewTableModel(columns, opts...), nil
}

func NewListCommand() *cobra.Command {
	config := NewListControllerConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
