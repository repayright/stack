package stack

import (
	"flag"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"

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

type ShowStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

func NewShowControllerConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack name")

	return config.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"show",
			"sh",
		},
		flags,
		config.Organization,
	)
}

type ShowController struct {
	store      *ShowStore
	config     *config.ControllerConfig
	fctlConfig *config.Config

	keymap *config.KeyMapHandler
}

func NewShowController(conf *config.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: conf,
		keymap: config.NewKeyMapHandler(),
	}
}

func (c *ShowController) GetStore() any {
	return c.store
}

func (c *ShowController) GetKeyMapAction() *config.KeyMapHandler {
	return c.keymap
}
func (c *ShowController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (config.Renderer, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := config.GetConfig(flags)

	if err != nil {
		return nil, err
	}
	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	var stack *membershipclient.Stack
	if len(c.config.GetArgs()) == 1 {
		if config.GetString(flags, internal.StackNameFlag) != "" {
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
		if config.GetString(flags, internal.StackNameFlag) == "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}
		stacksResponse, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacksResponse.Data {
			if s.Name == config.GetString(flags, internal.StackNameFlag) {
				stack = &s
				break
			}
		}
	}

	if stack == nil {
		return nil, errStackNotFound
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
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

func (c *ShowController) Render() (tea.Model, error) {
	model, err := internal.PrintStackInformation(c.config.GetOut(), c.config.GetAllFLags(), config.GetCurrentProfile(c.config.GetAllFLags(), c.fctlConfig), c.store.Stack, c.store.Versions)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func NewShowCommand() *cobra.Command {
	config := NewShowControllerConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithController(NewShowController(config)),
	)
}
