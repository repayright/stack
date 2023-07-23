package version

import (
	"flag"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	useVersion   = "version"
	shortVersion = "Get version"
)

type Store struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	Commit    string `json:"commit"`
}
type Controller struct {
	store  *Store
	config *config.ControllerConfig
}

var _ config.Controller[*Store] = (*Controller)(nil)

func NewStore() *Store {
	return &Store{
		Version:   fctl.Version,
		BuildDate: fctl.BuildDate,
		Commit:    fctl.Commit,
	}
}
func NewVersionConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useVersion, flag.ExitOnError)
	return config.NewControllerConfig(
		useVersion,
		shortVersion,
		shortVersion,
		[]string{"v"},
		flags,
	)
}
func NewController(config *config.ControllerConfig) *Controller {
	return &Controller{
		store:  NewStore(),
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
	return c, nil
}

func (c *Controller) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Version"), c.store.Version})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), c.store.BuildDate})
	tableData = append(tableData, []string{pterm.LightCyan("Commit"), c.store.Commit})
	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}
func NewCommand() *cobra.Command {
	c := NewVersionConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*Store](NewController(c)),
	)
}
