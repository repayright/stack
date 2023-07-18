package version

import (
	"flag"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	useVersion         = "version"
	descriptionVersion = "Get version"
	Version            = "develop"
	Commit             = "-"
	BuildDate          = "-"
)

type VersionStore struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	Commit    string `json:"commit"`
}
type VersionController struct {
	store  *VersionStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*VersionStore] = (*VersionController)(nil)

func NewDefaultVersionStore() *VersionStore {
	return &VersionStore{
		Version:   Version,
		BuildDate: BuildDate,
		Commit:    Commit,
	}
}
func NewVersionConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useVersion, flag.ExitOnError)
	c := fctl.NewControllerConfig(
		useVersion,
		descriptionVersion,
		[]string{"v"},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(descriptionVersion)

	return c
}
func NewVersionController(config fctl.ControllerConfig) *VersionController {
	return &VersionController{
		store:  NewDefaultVersionStore(),
		config: config,
	}
}

func (c *VersionController) GetStore() *VersionStore {
	return c.store
}

func (c *VersionController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *VersionController) Run() (fctl.Renderable, error) {
	return c, nil
}

func (c *VersionController) Render() error {
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
		fctl.WithShortDescription(*c.GetShortDescription()),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*VersionStore](NewVersionController(*c)),
	)
}
