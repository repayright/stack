package stack

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
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

func NewDefaultStackListStore() *StackListStore {
	return &StackListStore{
		Stacks: []Stack{},
	}
}

type StackListController struct {
	contextCommand *cobra.Command
	store          *StackListStore
	profile        *fctl.Profile
	organization   string
	config         *StackListControllerConfig
}

var _ modelutils.Controller[*StackListStore] = (*StackListController)(nil)

func NewStackListControllerConfig() *StackListControllerConfig {
	return &StackListControllerConfig{
		deletedFlag: false,
	}
}

func NewDefaultStackListControllerConfig() *StackListControllerConfig {

	return &StackListControllerConfig{
		deletedFlag: false,
	}
}

func (c *StackListController) GetContextualConfig() *StackListControllerConfig {
	return c.config
}

func NewStackListController() *StackListController {
	return &StackListController{
		store:  NewDefaultStackListStore(),
		config: New,
	}
}

func NewListCommand() *cobra.Command {

	c := NewStackListController()

	return fctl.NewMembershipCommand("list", fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List stacks"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithBoolFlag(deletedFlag, false, "Display deleted stacks"),
		fctl.WithController[*StackListStore](),
	)

}

func (c *StackListController) Init() {

	if c.config == nil {
		c.config = NewDefaultStackListControllerConfig()
	}

	// c.contextCommand = fctl.NewMembershipCommand(c.config.use, c.config.options...)

}

func (c *StackListController) GetStore() *StackListStore {
	return c.store
}

func (c *StackListController) Run() error {
	cmd := c.GetContextualCmd()
	if cmd == nil {
		return errors.New("contextual command is nil")
	}

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(cmd.Context(), organization).
		Deleted(fctl.GetBool(cmd, deletedFlag)).
		Execute()
	if err != nil {
		return errors.Wrap(err, "listing stacks")
	}

	c.profile = profile
	if len(rsp.Data) == 0 {
		return nil
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

	return nil
}

func (c *StackListController) Render() (modelutils.Model, error) {
	cmd := c.GetContextualCmd()
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

func (c *StackListController) GetKeyMapAction() *modelutils.KeyMapHandler[*StackListStore] {
	k := modelutils.NewKeyMapHandler[*StackListStore]()
	k.AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "Quit the application"),
		),
		func() *modelutils.Controller[*StackListStore] {
			return nil
		},
	)
	k.AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "move up"),
		),
		func() *modelutils.Controller[*StackListStore] {
			return nil
		},
	)
	k.AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "move down"),
		),
		func() *modelutils.Controller[*StackListStore] {
			return nil
		},
	)
	k.AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("? ", "Toggle help"),
		),
		func() *modelutils.Controller[*StackListStore] {
			return nil
		},
	)
	k.AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "show selected item"),
		),
		func() *modelutils.Controller[*StackListStore] {
			return nil
		},
	)

	return k
}
