package payments

import (
	"github.com/formancehq/fctl/cmd/payments/connectors"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("payments",
		fctl.WithShortDescription("Payments management"),
		fctl.WithChildCommands(
			connectors.NewConnectorsCommand(),
			NewListPaymentsCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
