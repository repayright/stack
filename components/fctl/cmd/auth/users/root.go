package users

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("users",
		fctl.WithShortDescription("Users management"),
		fctl.WithAliases("u", "user"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
