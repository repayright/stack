package accounts

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("accounts",
		fctl.WithAliases("acc", "a", "ac", "account"),
		fctl.WithShortDescription("Accounts management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewSetMetadataCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack, config.Ledger),
	)
}
