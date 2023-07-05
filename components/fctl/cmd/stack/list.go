package stack

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
)

// Every column in the table is a column in the struct
// The column name is suffixed with the max length of the column
// This is used to align the columns in the table
// e.g. maxLengthOrganizationId = 20
//
//	maxLengthStackId = 4
//
// Where the max length is the number of characters in the column name and values
const (
	maxLengthOrganizationId = 15
	maxLengthStackId        = 8
	maxLengthStackName      = 6
	maxLengthApiUrl         = 48
	maxLengthStackRegion    = 21
	maxLengthStackCreatedAt = 20
	maxLengthStackDeletedAt = 20
)

type Stack struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Dashboard string  `json:"dashboard"`
	RegionID  string  `json:"region"`
	CreatedAt string  `json:"created_at"`
	DeletedAt *string `json:"deleted_at"`
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
			CreatedAt: stack.CreatedAt.Format(time.RFC3339),
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

func (c *StackListController) Render(cmd *cobra.Command, args []string) error {
	// Default Columns
	columns := ui.NewArrayColumn(
		ui.NewColumn("Organization Id", maxLengthOrganizationId),
		ui.NewColumn("Stack Id", maxLengthStackId),
		ui.NewColumn("Name", maxLengthStackName),
		ui.NewColumn("API URL", maxLengthApiUrl),
		ui.NewColumn("Region", maxLengthStackRegion),
		ui.NewColumn("Created At", maxLengthStackCreatedAt),
	)

	// Add Deleted At column if --deleted flag is set
	if fctl.GetBool(cmd, deletedFlag) {
		columns = columns.AddColumn("Deleted At", maxLengthStackDeletedAt)
	}

	// Create table data
	tableData := fctl.Map(c.store.Stacks, func(stack Stack) table.Row {
		data := []string{
			c.organization,
			stack.Id,
			stack.Name,
			stack.Dashboard,
			stack.RegionID,
			stack.CreatedAt,
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

	// Default table options
	opts := ui.NewTableOptions(columns, tableData)

	// Add plain table option if --plain flag is set
	isPlain := fctl.GetString(cmd, fctl.OutputFlag) == "plain"

	var opt table.Option
	if isPlain {
		opt = ui.WithHeight(len(tableData))
	} else {
		opt = ui.WithHeight(ui.MaxTableHeight)
	}

	opts = append(opts, opt)

	t := ui.NewTableModel(opts...)

	return displayStackList(cmd, t)
}

func displayStackList(cmd *cobra.Command, t *ui.TableModel) error {
	isPlain := fctl.GetString(cmd, fctl.OutputFlag)

	if isPlain == "plain" {
		fmt.Println(t.View())
		return nil
	}

	header := ui.NewHeader()

	d := ui.NewDisplay().AppendModels(t).SetHeader(header)

	if _, err := tea.NewProgram(d, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}
