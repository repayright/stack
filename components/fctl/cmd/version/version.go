package version

import (
	"flag"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	useVersion   = "version"
	shortVersion = "Get version"
	Version      = "develop"
	Commit       = "-"
	BuildDate    = "-"
)

type Store struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	Commit    string `json:"commit"`
}
type Controller struct {
	store  *Store
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*Store] = (*Controller)(nil)

func NewDefaultStore() *Store {
	return &Store{
		Version:   Version,
		BuildDate: BuildDate,
		Commit:    Commit,
	}
}
func NewVersionConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useVersion, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useVersion,
		shortVersion,
		shortVersion,
		[]string{"v"},
		os.Stdout,
		flags,
	)
}
func NewController(config *fctl.ControllerConfig) *Controller {
	return &Controller{
		store:  NewDefaultStore(),
		config: config,
	}
}

func (c *Controller) GetStore() *Store {
	return c.store
}

func (c *Controller) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *Controller) Run() (fctl.Renderable, error) {
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
