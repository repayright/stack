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
	useDelete              = "delete <name>"
	shortDescriptionDelete = "Delete a profile"
	descriptionDelete      = "Delete a profile"
)

type DeleteStore struct {
	Success bool `json:"success"`
}
type DeleteController struct {
	store  *DeleteStore
	config fctl.ControllerConfig
}

func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useDelete,
		descriptionDelete,
		[]string{
			"del",
		},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortDescriptionDelete)

	return c
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{
		Success: false,
	}
}

func NewDeleteController(config fctl.ControllerConfig) *DeleteController {
	return &DeleteController{
		store:  NewDeleteStore(),
		config: config,
	}
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *DeleteController) Run() (fctl.Renderable, error) {

	flags := c.config.GetFlags()
	args := c.config.GetArgs()

	if len(args) < 1 {
		return nil, errors.New("Profile(delete): invalid number of arguments")
	}

	config, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	if err := config.DeleteProfile(args[0]); err != nil {
		return nil, err
	}

	if err := config.Persist(); err != nil {
		return nil, errors.Wrap(err, "updating config")
	}

	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Profile deleted!")
	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription(*config.GetShortDescription()),
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController[*DeleteStore](NewDeleteController(*config)),
	)
}
