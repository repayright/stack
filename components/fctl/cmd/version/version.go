package version

import (
	"fmt"

	blist "github.com/charmbracelet/bubbles/list"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/list"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type VersionStore struct {
	Version   string `json:"version" yaml:"version"`
	BuildDate string `json:"buildDate" yaml:"buildDate"`
	Commit    string `json:"commit" yaml:"commit"`
}
type VersionController struct {
	store *VersionStore
}

var _ fctl.Controller[*VersionStore] = (*VersionController)(nil)

func NewDefaultVersionStore() *VersionStore {
	return &VersionStore{
		Version:   fctl.Version,
		BuildDate: "-",
		Commit:    "-",
	}
}

func NewVersionController() *VersionController {
	return &VersionController{
		store: NewDefaultVersionStore(),
	}
}

func NewCommand() *cobra.Command {
	return fctl.NewCommand("version",
		fctl.WithShortDescription("Get version"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*VersionStore](NewVersionController()),
	)
}

func (c *VersionController) GetStore() *VersionStore {
	return c.store
}

func (c *VersionController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *VersionController) Render(cmd *cobra.Command, args []string) (ui.Model, error) {

	// Default List
	items := []blist.Item{
		list.NewItem(pterm.LightCyan("Version"), c.store.Version),
		list.NewItem(pterm.LightCyan("Date"), c.store.BuildDate),
		list.NewItem(pterm.LightCyan("Commit"), c.store.Commit),
	}

	model, err := list.NewDefaultListModel(items, false)
	if err != nil {
		return nil, err
	}
	model = model.WithTitle("FCTL Information")

	//Print the list
	fmt.Println(model.View())
	return nil, nil

}
