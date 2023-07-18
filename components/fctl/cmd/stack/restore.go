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
	useRestore         = "restore <stack-id>"
	descriptionRestore = "Restore a stack"
)

type StackRestoreStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewDefaultVersionStore() *StackRestoreStore {
	return &StackRestoreStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

func NewStackRestoreControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useRestore, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack name")

	return fctl.NewControllerConfig(
		useRestore,
		descriptionRestore,
		[]string{
			"restore",
			"re",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*StackRestoreStore] = (*StackRestoreController)(nil)

type StackRestoreController struct {
	store      *StackRestoreStore
	config     fctl.ControllerConfig
	fctlConfig *fctl.Config
}

func NewStackRestoreController(config fctl.ControllerConfig) *StackRestoreController {
	return &StackRestoreController{
		store:  NewDefaultVersionStore(),
		config: config,
	}
}

func (c *StackRestoreController) GetStore() *StackRestoreStore {
	return c.store
}

func (c *StackRestoreController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *StackRestoreController) Run() (fctl.Renderable, error) {
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

	if len(c.config.GetArgs()) == 0 {
		return nil, fmt.Errorf("stack id is required")
	}

	response, _, err := apiClient.DefaultApi.
		RestoreStack(ctx, organization, c.config.GetArgs()[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	if err := waitStackReady(ctx, flags, profile, response.Data); err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, response.Data)
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

	c.store.Stack = response.Data
	c.store.Versions = versions.GetVersionsResponse
	c.fctlConfig = cfg

	return c, nil
}

func (c *StackRestoreController) Render() error {
	return internal.PrintStackInformation(c.config.GetOut(), fctl.GetCurrentProfile(c.config.GetAllFLags(), c.fctlConfig), c.store.Stack, c.store.Versions)
}

func NewRestoreStackCommand() *cobra.Command {
	config := NewStackRestoreControllerConfig()
	return fctl.NewMembershipCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*StackRestoreStore](NewStackRestoreController(*config)),
	)
}
