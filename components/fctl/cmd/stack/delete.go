package stack

import (
	"flag"
	"os"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete         = "delete (<stack-id> | --name=<stack-name>)"
	descriptionDelete = "Delete a stack"
)

type DeletedStackStore struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}

func NewDefaultDeletedStackStore() *DeletedStackStore {
	return &DeletedStackStore{
		Stack:  &membershipclient.Stack{},
		Status: "",
	}
}

func NewStackDeleteControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack to remove")
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useDelete,
		descriptionDelete,
		[]string{
			"delete",
			"del",
			"rm",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*DeletedStackStore] = (*StackDeleteController)(nil)

type StackDeleteController struct {
	store  *DeletedStackStore
	config fctl.ControllerConfig
}

func NewStackDeleteController(config fctl.ControllerConfig) *StackDeleteController {
	return &StackDeleteController{
		store:  NewDefaultDeletedStackStore(),
		config: config,
	}
}

func (c *StackDeleteController) GetStore() *DeletedStackStore {
	return c.store
}

func (c *StackDeleteController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *StackDeleteController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	var stack *membershipclient.Stack
	if len(c.config.GetArgs()) == 1 {
		if fctl.GetString(flags, internal.StackNameFlag) != "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}

		rsp, _, err := apiClient.DefaultApi.ReadStack(ctx, organization, c.config.GetArgs()[0]).Execute()
		if err != nil {
			return nil, err
		}
		stack = rsp.Data
	} else {
		if fctl.GetString(flags, internal.StackNameFlag) == "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}
		stacks, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacks.Data {
			if s.Name == fctl.GetString(flags, internal.StackNameFlag) {
				stack = &s
				break
			}
		}
	}
	if stack == nil {
		return nil, errors.New("Stack not found")
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to delete stack '%s'", stack.Name) {
		return nil, fctl.ErrMissingApproval
	}

	if _, err := apiClient.DefaultApi.DeleteStack(ctx, organization, stack.Id).Execute(); err != nil {
		return nil, errors.Wrap(err, "deleting stack")
	}

	c.store.Stack = stack
	c.store.Status = "OK"

	return c, nil
}

func (c *StackDeleteController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Stack deleted.")
	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewStackDeleteControllerConfig()
	return fctl.NewMembershipCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithAliases(config.GetAliases()...),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithController[*DeletedStackStore](NewStackDeleteController(*config)),
	)
}
