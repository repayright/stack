package invitations

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("invitations",
		fctl.WithAliases("invit", "inv", "i"),
		fctl.WithShortDescription("Invitations management"),
		fctl.WithChildCommands(
			NewSendCommand(),
			NewListCommand(),
			NewDeleteCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
