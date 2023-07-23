package invitations

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("invitations",
		fctl.WithShortDescription("Invitations management"),
		fctl.WithAliases("invit", "inv", "i"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewAcceptCommand(),
			NewDeclineCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
