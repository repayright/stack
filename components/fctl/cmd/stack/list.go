package stack

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/helpers"
	"github.com/formancehq/fctl/pkg/modelutils"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	uitable "github.com/formancehq/fctl/pkg/components/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// This defines the minimum length of the columns in the table
const (
	minLengthOrganizationId = 15
	minLengthStackId        = 8
	minLengthStackName      = 10
	minLengthApiUrl         = 49
	minLengthStackRegion    = 36
	minLengthStackCreatedAt = 21
	minLengthStackDeletedAt = 21
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

func NewListControllerConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.Bool(deletedFlag, false, "Show deleted stacks")

	return config.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"list",
			"ls",
		},
		flags,
		config.Organization,
	)
}

var _ config.Controller = (*ListController)(nil)

type ListController struct {
	store        *ListStore
	profile      *config.Profile
	config       *config.ControllerConfig
	organization string
	keyMapAction *config.KeyMapHandler
}

func NewListController(config *config.ControllerConfig) *ListController {
	return &ListController{
		store:        NewListStore(),
		config:       config,
		keyMapAction: NewKeyMapAction(),
	}
}

func (c *ListController) GetStore() any {
	return c.store
}

func (c *ListController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (config.Renderer, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := config.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := config.GetCurrentProfile(flags, cfg)

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).
		Deleted(config.GetBool(flags, deletedFlag)).
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

func (c *ListController) Render() (tea.Model, error) {

	flags := c.config.GetAllFLags()

	rows := fctl.Map(c.store.Stacks, func(stack Stack) *uitable.Row {
		data := []string{
			c.organization,
			stack.Id,
			stack.Name,
			stack.Dashboard,
			stack.RegionID,
			*stack.CreatedAt,
		}

		if config.GetBool(flags, deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, *stack.DeletedAt)
			} else {
				data = append(data, "")
			}
		}

		cells := fctl.Map(data, func(s string) *uitable.Cell {
			return uitable.NewCell(s)
		})

		return uitable.NewRow(cells...)
	})

	row := uitable.NewRow(
		uitable.NewCell("Organization Id", uitable.WithWidth(minLengthOrganizationId)),
		uitable.NewCell("Stack Id", uitable.WithWidth(minLengthStackId)),
		uitable.NewCell("Name", uitable.WithWidth(minLengthStackName)),
		uitable.NewCell("API URL", uitable.WithWidth(minLengthApiUrl)),
		uitable.NewCell("Region", uitable.WithWidth(minLengthStackRegion)),
		uitable.NewCell("Created At", uitable.WithWidth(minLengthStackCreatedAt)),
	)

	if config.GetBool(flags, deletedFlag) {
		row.AddCell(uitable.NewCell("Deleted At", uitable.WithWidth(minLengthStackDeletedAt)))
	}

	isPlain := config.GetString(flags, config.OutputFlag) == "plain"
	if isPlain {
		return uitable.NewTable(
			row,
			rows,
			uitable.WithDefaultStyle(),
			uitable.WithFullScreen(false),
		), nil
	}

	return uitable.NewTable(
		row,
		rows,
		uitable.WithDefaultStyle(),
		uitable.WithFullScreen(true),
	), nil
}

func NewKeyMapAction() *config.KeyMapHandler {
	return config.NewKeyMapHandler().AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "move up"),
		),
		func(m tea.Model) tea.Msg {
			return nil
		},
	).AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "move down"),
		),
		func(m tea.Model) tea.Msg {
			return nil
		},
	).AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("? ", "Toggle help"),
		),
		func(m tea.Model) tea.Msg {
			return nil
		},
	).AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "show selected item"),
		),
		func(m tea.Model) tea.Msg {
			//Cast model to table.Model
			t, ok := m.(*uitable.Table)
			if !ok {
				panic("invalid model type")
				return nil
			}

			selectedRow := t.SelectedRow()
			id := selectedRow.Items()[1].String()
			log := helpers.NewLogger("SELECTED")
			log.Log("ID", id)
			c := NewShowControllerConfig()
			controller := NewShowController(c)
			c.SetOut(os.Stdout)
			c.SetContext(context.TODO())
			c.SetArgs([]string{id})

			return modelutils.ChangeViewMsg{
				Controller: controller,
			}
		},
	).AddNewKeyBinding(
		key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "delete selected item"),
		),
		func(m tea.Model) tea.Msg {
			//Cast model to table.Model
			t, ok := m.(*uitable.Table)
			if !ok {
				panic("invalid model type")
				return nil
			}

			selectedRow := t.SelectedRow()
			id := selectedRow.Items()[1].String()

			c := NewShowControllerConfig()
			controller := NewShowController(c)
			c.SetOut(os.Stdout)
			c.SetContext(context.TODO())
			c.SetArgs([]string{id})

			return modelutils.ChangeViewMsg{
				Controller: controller,
			}
		},
	)
}

func (c *ListController) GetKeyMapAction() *config.KeyMapHandler {
	return c.keyMapAction
}

func NewListCommand() *cobra.Command {
	c := NewListControllerConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController(NewListController(c)),
	)
}
