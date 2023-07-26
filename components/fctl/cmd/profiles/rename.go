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
	useRename   = "rename <old-name> <new-name>"
	shortRename = "Rename a profile"
)

type RenameStore struct {
	Success bool `json:"success"`
}

func NewRenameStore() *RenameStore {
	return &RenameStore{
		Success: false,
	}
}
func NewRenameConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useRename, flag.ExitOnError)

	return config.NewControllerConfig(
		useRename,
		shortRename,
		shortRename,
		[]string{},
		flags,
	)
}

var _ config.Controller = (*RenameController)(nil)

type RenameController struct {
	store  *RenameStore
	config *config.ControllerConfig
}

func NewRenameController(config *config.ControllerConfig) *RenameController {
	return &RenameController{
		store:  NewRenameStore(),
		config: config,
	}
}

func (c *RenameController) GetStore() any {
	return c.store
}

func (c *RenameController) GetKeyMapAction() *config.KeyMapHandler {
	return nil
}
func (c *RenameController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *RenameController) Run() (config.Renderer, error) {
	flags := c.config.GetFlags()
	args := c.config.GetArgs()

	oldName := args[0]
	newName := args[1]

	config, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	p := config.GetProfile(oldName)
	if p == nil {
		return nil, errors.New("profile not found")
	}

	if err := config.DeleteProfile(oldName); err != nil {
		return nil, err
	}
	if config.GetCurrentProfileName() == oldName {
		config.SetCurrentProfile(newName, p)
	} else {
		config.SetProfile(newName, p)
	}

	if err := config.Persist(); err != nil {
		return nil, errors.Wrap(config.Persist(), "Updating config")
	}

	c.store.Success = true
	return c, nil
}

func (c *RenameController) Render() (tea.Model, error) {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Profile renamed!")
	return nil, nil
}

func NewRenameCommand() *cobra.Command {
	config := NewRenameConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController(NewRenameController(config)),
	)
}
