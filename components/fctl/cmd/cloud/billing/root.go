package billing

import (
	fctl "github.com/formancehq/fctl/pkg"
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
		fctl.WithScopesFlags(fctl.Organization, fctl.Stack),
	)
}
