package me

import (
	"errors"
	"flag"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useInfo   = "info"
	shortInfo = "Display user information"
)

type InfoStore struct {
	Subject string `json:"subject"`
	Email   string `json:"email"`
}

func NewInfoStore() *InfoStore {
	return &InfoStore{}
}

func NewInfoConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useInfo, flag.ExitOnError)
	return config.NewControllerConfig(
		useInfo,
		shortInfo,
		shortInfo,
		[]string{
			"i", "in",
		},
		flags,
	)
}

var _ config.Controller[*InfoStore] = (*InfoController)(nil)

type InfoController struct {
	store  *InfoStore
	config *config.ControllerConfig
}

func NewInfoController(config *config.ControllerConfig) *InfoController {
	return &InfoController{
		store:  NewInfoStore(),
		config: config,
	}
}

func (c *InfoController) GetStore() *InfoStore {
	return c.store
}

func (c *InfoController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *InfoController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	if !profile.IsConnected() {
		return nil, errors.New("Not logged. Use 'login' command before.")
	}

	userInfo, err := profile.GetUserInfo()
	if err != nil {
		return nil, err
	}

	c.store.Subject = userInfo.Subject
	c.store.Email = userInfo.Email

	return c, nil
}

func (c *InfoController) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Subject"), c.store.Subject})
	tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.Email})

	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewInfoCommand() *cobra.Command {
	config := NewInfoConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*InfoStore](NewInfoController(config)),
	)
}
