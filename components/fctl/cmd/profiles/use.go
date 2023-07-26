package profiles

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/cmd/profiles/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useProfile         = "use <name>"
	shortProfile       = "Use a profile"
	descriptionProfile = "Select a profile to use"
)

type UseStore struct {
	Success bool `json:"success"`
}

func NewUseStore() *UseStore {
	return &UseStore{
		Success: false,
	}
}

func NewUseConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useProfile, flag.ExitOnError)

	return config.NewControllerConfig(
		useProfile,
		descriptionProfile,
		shortProfile,
		[]string{
			"u",
		},
		flags,
	)
}

type UseController struct {
	store  *UseStore
	config *config.ControllerConfig
}

var _ config.Controller = (*UseController)(nil)

func NewUseController(config *config.ControllerConfig) *UseController {
	return &UseController{
		store:  NewUseStore(),
		config: config,
	}
}

func (c *UseController) GetKeyMapAction() *config.KeyMapHandler {
	return nil
}

func (c *UseController) GetStore() any {
	return c.store
}

func (c *UseController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *UseController) Run() (config.Renderer, error) {
	flags := c.config.GetAllFLags()
	args := c.config.GetArgs()

	if len(args) < 1 {
		return nil, errors.New("No profile name provided")
	}

	config, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	config.SetCurrentProfileName(args[0])
	if err := config.Persist(); err != nil {
		return nil, errors.Wrap(err, "Updating config")
	}

	c.store.Success = true
	return c, nil
}

func (c *UseController) Render() (tea.Model, error) {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Selected profile updated!")
	return nil, nil
}

func NewUseCommand() *cobra.Command {
	config := NewUseConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController(NewUseController(config)),
	)
}
