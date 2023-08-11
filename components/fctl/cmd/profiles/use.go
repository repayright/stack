package profiles

import (
	"flag"
	"github.com/formancehq/fctl/cmd/profiles/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

const (
	useProfile         = "use <name>"
	shortDescription   = "Use a profile"
	descriptionProfile = "Select a profile to use"
)

type UseStore struct {
	Success bool `json:"success"`
}

func NewDefaultUseStore() *UseStore {
	return &UseStore{
		Success: false,
	}
}

func NewUseConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useProfile, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useProfile,
		descriptionProfile,
		[]string{
			"u",
		},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortDescription)

	return c
}

type UseController struct {
	store  *UseStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*UseStore] = (*UseController)(nil)

func NewUseController(config fctl.ControllerConfig) *UseController {
	return &UseController{
		store:  NewDefaultUseStore(),
		config: config,
	}
}

func (c *UseController) GetStore() *UseStore {
	return c.store
}

func (c *UseController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *UseController) Run() (fctl.Renderable, error) {
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

func (c *UseController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Selected profile updated!")
	return nil
}

func NewUseCommand() *cobra.Command {
	config := NewUseConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithAliases(config.GetAliases()...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription(*config.GetShortDescription()),
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController[*UseStore](NewUseController(*config)),
	)
}
