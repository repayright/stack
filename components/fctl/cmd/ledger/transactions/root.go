package transactions

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("transactions",
		fctl.WithAliases("t", "txs", "tx"),
		fctl.WithShortDescription("Transactions management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewNumCommand(),
			NewRevertCommand(),
			NewShowCommand(),
			NewSetMetadataCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack, config.Ledger),
	)
}
