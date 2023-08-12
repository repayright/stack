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
	useRestore   = "restore <stack-id>"
	shortRestore = "Restore a stack"
)

type RestoreStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewRestoreStore() *RestoreStore {
	return &RestoreStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

func NewRestoreConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useRestore, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack name")

	return config.NewControllerConfig(
		useRestore,
		shortRestore,
		shortRestore,
		[]string{
			"restore",
			"re",
		},
		flags,
		config.Organization,
	)
}

//var _ config.Controller[*RestoreStore] = (*RestoreController)(nil)

type RestoreController struct {
	store      *RestoreStore
	config     *config.ControllerConfig
	fctlConfig *config.Config
}

func (c *RestoreController) GetKeyMapAction() *config.KeyMapHandler {
	return nil
}

func NewRestoreController(config *config.ControllerConfig) *RestoreController {
	return &RestoreController{
		store:  NewRestoreStore(),
		config: config,
	}
}

func (c *RestoreController) GetStore() any {
	return c.store
}

func (c *RestoreController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *RestoreController) Run() (config.Renderer, error) {
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

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
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

	profile := config.GetCurrentProfile(flags, cfg)

	if err := waitStackReady(ctx, c.config.GetOut(), flags, profile, response.Data); err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, response.Data, out)
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

func (c *RestoreController) Render() (tea.Model, error) {
	model, err := internal.PrintStackInformation(c.config.GetOut(), c.config.GetAllFLags(), config.GetCurrentProfile(c.config.GetAllFLags(), c.fctlConfig), c.store.Stack, c.store.Versions)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func NewRestoreStackCommand() *cobra.Command {
	config := NewRestoreConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController(NewRestoreController(config)),
	)
}
