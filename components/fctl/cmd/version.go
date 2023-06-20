package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	Version = "develop"
)

type VersionStruct struct {
	Version   string `json:"version" yaml:"version"`
	BuildDate string `json:"buildDate" yaml:"buildDate"`
	Commit    string `json:"commit" yaml:"commit"`
}
type VersionController struct {
	store *fctl.SharedStore
}

func NewVersion() *VersionController {
	return &VersionController{
		store: fctl.NewSharedStore(),
	}
}

func NewVersionCommand() *cobra.Command {
	return fctl.NewCommand("version",
		fctl.WithShortDescription("Get version"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController(NewVersion()),
	)
}

func (c *VersionController) GetStore() *fctl.SharedStore {
	return c.store
}

func (c *VersionController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	version := &VersionStruct{
		Version:   "develop",
		BuildDate: "-",
		Commit:    "-",
	}

	c.GetStore().SetData(version)

	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *VersionController) Render(cmd *cobra.Command, args []string) error {
	data := c.GetStore().GetData().(*VersionStruct)

	// Default List
	items := []list.Item{
		ui.NewItem(pterm.LightCyan("Version"), data.Version),
		ui.NewItem(pterm.LightCyan("Date"), data.BuildDate),
		ui.NewItem(pterm.LightCyan("Commit"), data.Commit),
	}

	model, err := ui.NewDefaultListModel(items, false)
	if err != nil {
		return err
	}
	model = model.WithTitle("FCTL Information")

	//Print the list
	fmt.Println(model.View())
	return nil

}
