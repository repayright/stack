package ui

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

const (
	useUI         = "ui"
	shortUI       = "Open UI"
	descriptionUI = "Open UI in browser (if available), otherwise print the url to the console."
)

type Store struct {
	UIUrl        string `json:"stackUrl"`
	FoundBrowser bool   `json:"browserFound"`
}

type Controller struct {
	store  *Store
	config *config.ControllerConfig
}

var _ config.Controller[*Store] = (*Controller)(nil)

func NewDefaultUiStore() *Store {
	return &Store{
		UIUrl:        "",
		FoundBrowser: false,
	}
}
func NewUiConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useUI, flag.ExitOnError)

	return config.NewControllerConfig(
		useUI,
		descriptionUI,
		shortUI,
		[]string{},
		flags,
		config.Organization, config.Stack,
	)
}

func NewController(config *config.ControllerConfig) *Controller {
	return &Controller{
		store:  NewDefaultUiStore(),
		config: config,
	}
}

func (c *Controller) GetStore() *Store {
	return c.store
}

func (c *Controller) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *Controller) Run() (modelutils.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organization, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	stackUrl := profile.ServicesBaseUrl(stack)

	c.store.UIUrl = stackUrl.String()

	if err := OpenUrl(c.store.UIUrl); err != nil {
		c.store.FoundBrowser = true
	}

	return c, nil
}

func (c *Controller) Render() error {

	fmt.Fprintln(c.config.GetOut(), "Opening url: ", c.store.UIUrl)

	return nil
}

func NewCommand() *cobra.Command {
	config := NewUiConfig()
	return fctl.NewCommand(useUI,
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*Store](NewController(config)),
	)
}
