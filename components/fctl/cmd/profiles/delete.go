package profiles

import (
	"flag"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	"github.com/formancehq/fctl/cmd/profiles/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete         = "delete <name>"
	shortDelete       = "Delete a profile"
	descriptionDelete = "Delete a profile"
)

type DeleteStore struct {
	Success bool `json:"success"`
}

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{
		Success: false,
	}
}
func NewDeleteConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)

	return config.NewControllerConfig(
		useDelete,
		descriptionDelete,
		shortDelete,
		[]string{
			"del",
		},
		flags,
	)

}

var _ config.Controller[*DeleteStore] = (*DeleteController)(nil)

type DeleteController struct {
	store  *DeleteStore
	config *config.ControllerConfig
}

func NewDeleteController(config *config.ControllerConfig) *DeleteController {
	return &DeleteController{
		store:  NewDeleteStore(),
		config: config,
	}
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *DeleteController) Run() (modelutils.Renderable, error) {

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
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController[*DeleteStore](NewDeleteController(config)),
	)
}
