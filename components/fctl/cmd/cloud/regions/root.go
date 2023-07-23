package regions

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("regions",
		fctl.WithAliases("region", "reg"),
		fctl.WithShortDescription("Regions management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
