package billing

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("billing",
		fctl.WithAliases("bil", "b"),
		fctl.WithShortDescription("Billing management"),
		fctl.WithChildCommands(
			NewPortalCommand(),
			NewSetupCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
