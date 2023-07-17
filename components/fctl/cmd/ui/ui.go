package ui

import (
	"flag"
	"fmt"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

const (
	useUI              = "ui"
	shortDescriptionUI = "Open UI"
	descriptionUI      = "Open UI in browser (if available), otherwise print the url to the console."
)

type UiStruct struct {
	UIUrl        string `json:"stackUrl"`
	FoundBrowser bool   `json:"browserFound"`
}

type UiController struct {
	store  *UiStruct
	config fctl.ControllerConfig
}

var _ fctl.Controller[*UiStruct] = (*UiController)(nil)

func NewDefaultUiStore() *UiStruct {
	return &UiStruct{
		UIUrl:        "",
		FoundBrowser: false,
	}
}
func NewUiConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useUI, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useUI,
		descriptionUI,
		[]string{},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortDescriptionUI)

	return c
}

func NewUiController(config fctl.ControllerConfig) *UiController {
	return &UiController{
		store:  NewDefaultUiStore(),
		config: config,
	}
}

func (c *UiController) GetStore() *UiStruct {
	return c.store
}

func (c *UiController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *UiController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organization)
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

func (c *UiController) Render() error {

	fmt.Println("Opening url: ", c.store.UIUrl)

	return nil
}

func NewCommand() *cobra.Command {
	config := NewUiConfig()
	return fctl.NewStackCommand(useUI,
		fctl.WithShortDescription(*config.GetShortDescription()),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*UiStruct](NewUiController(*config)),
	)
}
