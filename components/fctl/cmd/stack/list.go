package stack

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
	useList     = "list"
	description = "List stacks"
)

type Stack struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Dashboard string  `json:"dashboard"`
	RegionID  string  `json:"region"`
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

func NewStackListControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.Bool(deletedFlag, false, "Show deleted stacks")

	return fctl.NewControllerConfig(
		useList,
		description,
		[]string{
			"list",
			"ls",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*StackListStore] = (*StackListController)(nil)

type StackListController struct {
	store   *StackListStore
	profile *fctl.Profile
	config  fctl.ControllerConfig
}

func NewStackListController(config fctl.ControllerConfig) *StackListController {
	return &StackListController{
		store:  NewDefaultStackListStore(),
		config: config,
	}
}

func (c *StackListController) GetStore() *StackListStore {
	return c.store
}

func (c *StackListController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *StackListController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg)
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

	c.store.Stacks = fctl.Map(rsp.Data, func(stack membershipclient.Stack) Stack {
		return Stack{
			Id:        stack.Id,
			Name:      stack.Name,
			Dashboard: c.profile.ServicesBaseUrl(&stack).String(),
			RegionID:  stack.RegionID,
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

func (c *StackListController) Render() error {
	if len(c.store.Stacks) == 0 {
		fmt.Fprintln(os.Stdout, "No stacks found.")
		return nil
	}

	tableData := fctl.Map(c.store.Stacks, func(stack Stack) []string {
		data := []string{
			stack.Id,
			stack.Name,
			stack.Dashboard,
			stack.RegionID,
		}
		if fctl.GetBool(c.config.GetAllFLags(), deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, *stack.DeletedAt)
			} else {
				data = append(data, "")
			}
		}
		return data
	})

	headers := []string{"ID", "Name", "Dashboard", "Region"}
	if fctl.GetBool(c.config.GetAllFLags(), deletedFlag) {
		headers = append(headers, "Deleted at")
	}

	tableData = fctl.Prepend(tableData, headers)

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(os.Stdout).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {
	config := NewStackListControllerConfig()

	return fctl.NewMembershipCommand(config.GetUse(), fctl.WithAliases(config.GetAliases()...),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithGoFlagSet(config.GetFlags()),
		fctl.WithController[*StackListStore](NewStackListController(*config)),
	)
}
