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

func NewVersionCommand() *cobra.Command {
	return fctl.NewCommand("version",
		fctl.WithShortDescription("Get version"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithRunE(versionCommand),
		fctl.WrapOutputPostRunE(view),
	)
}

func versionCommand(cmd *cobra.Command, args []string) error {

	version := &VersionStruct{
		Version:   "develop",
		BuildDate: "-",
		Commit:    "-",
	}

	fctl.SetSharedData(version, nil, nil, nil)

	return nil
}

// TODO: This need to use the ui.NewListModel
func view(cmd *cobra.Command, args []string) error {
	data := fctl.GetSharedData().(*VersionStruct)

	// Default List
	items := []list.Item{
		ui.NewItem(pterm.LightCyan("Version"), data.Version),
		ui.NewItem(pterm.LightCyan("Date"), data.BuildDate),
		ui.NewItem(pterm.LightCyan("Commit"), data.Commit),
	}

	model := ui.NewDefaultListModel(items, false).WithTitle("FCTL Information")

	//Print the list
	fmt.Println(model.View())
	return nil

}
