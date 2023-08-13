package profiles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewNodeController() *config.Node {
	return config.NewConfigNode(*config.NewControllerConfig(
		"profiles",
		"Profiles management",
		"profiles",
		[]string{"profiles", "p", "prof"},
		nil,
	),
		NewListController(NewListConfig()),
	)
}
func NewCommand() *cobra.Command {
	return fctl.NewCommand("profiles",
		fctl.WithAliases("p", "prof"),
		fctl.WithShortDescription("Profiles management"),
		fctl.WithChildCommands(
			NewDeleteCommand(),
			NewListCommand(),
			NewRenameCommand(),
			NewShowCommand(),
			NewUseCommand(),
			NewSetDefaultOrganizationCommand(),
		),
	)
}
