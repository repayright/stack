package stack

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	useShow   = "show (<stack-id> | --name=<stack-name>)"
	shortShow = "Show a stack"
)

var errStackNotFound = errors.New("stack not found")

type StackShowStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewDefaultStackShowStore() *StackShowStore {
	return &StackShowStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

func NewStackShowControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack name")

	return fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"show",
			"sh",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*StackShowStore] = (*StackShowController)(nil)

type StackShowController struct {
	store      *StackShowStore
	config     fctl.ControllerConfig
	fctlConfig *fctl.Config
}

func NewStackShowController(config fctl.ControllerConfig) *StackShowController {
	return &StackShowController{
		store:  NewDefaultStackShowStore(),
		config: config,
	}
}

func (c *StackShowController) GetStore() *StackShowStore {
	return c.store
}

func (c *StackShowController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *StackShowController) Run() (fctl.Renderable, error) {
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
		stackResponse, httpResponse, err := apiClient.DefaultApi.ReadStack(ctx, organization, c.config.GetArgs()[0]).Execute()
		if err != nil {
			if httpResponse.StatusCode == http.StatusNotFound {
				return nil, errStackNotFound
			}
			return nil, errors.Wrap(err, "listing stacks")
		}
		stack = stackResponse.Data
	} else {
		if fctl.GetString(flags, internal.StackNameFlag) == "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}
		stacksResponse, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacksResponse.Data {
			if s.Name == fctl.GetString(flags, internal.StackNameFlag) {
				stack = &s
				break
			}
		}
	}

	if stack == nil {
		return nil, errStackNotFound
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, err
	}

	versions, err := stackClient.GetVersions(ctx)
	if err != nil {
		return nil, err
	}

	if versions.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
	}

	c.store.Stack = stack
	c.store.Versions = versions.GetVersionsResponse
	c.fctlConfig = cfg

	return c, nil

}

func (c *StackShowController) Render() error {
	return internal.PrintStackInformation(c.config.GetOut(), fctl.GetCurrentProfile(c.config.GetAllFLags(), c.fctlConfig), c.store.Stack, c.store.Versions)
}

func NewShowCommand() *cobra.Command {
	config := NewStackShowControllerConfig()
	return fctl.NewMembershipCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithController[*StackShowStore](NewStackShowController(*config)),
	)
}
