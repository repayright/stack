package stack

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
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

type Stack struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Dashboard string  `json:"dashboard"`
	RegionID  string  `json:"region"`
	CreatedAt *string `json:"createdAt"`
	DeletedAt *string `json:"deletedAt"`
}
type StackListStore struct {
	Stacks []Stack `json:"stacks"`
}

type StackListController struct {
	store        *StackListStore
	profile      *fctl.Profile
	organization string
}

var _ fctl.Controller[*StackListStore] = (*StackListController)(nil)

func NewDefaultStackListStore() *StackListStore {
	return &StackListStore{
		Stacks: []Stack{},
	}
}

func NewStackListController() *StackListController {
	return &StackListController{
		store: NewDefaultStackListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewMembershipCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List stacks"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithBoolFlag(deletedFlag, false, "Display deleted stacks"),
		fctl.WithController[*StackListStore](NewStackListController()),
	)
}
func (c *StackListController) GetStore() *StackListStore {
	return c.store
}

func (c *StackListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(cmd.Context(), organization).
		Deleted(fctl.GetBool(cmd, deletedFlag)).
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

func (c *StackListController) Render(cmd *cobra.Command, args []string) (ui.Model, error) {

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

		if fctl.GetBool(cmd, deletedFlag) {
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
	isPlain := fctl.GetString(cmd, fctl.OutputFlag) == "plain"
	// Default Columns
	columns = ui.NewArrayColumn(
		ui.NewColumn("Organization Id", minLengthOrganizationId),
		ui.NewColumn("Stack Id", minLengthStackId),
		ui.NewColumn("Name", minLengthStackName),
		ui.NewColumn("API URL", minLengthApiUrl),
		ui.NewColumn("Region", minLengthStackRegion),
		ui.NewColumn("Created At", minLengthStackCreatedAt),
	)
	if fctl.GetBool(cmd, deletedFlag) {
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
