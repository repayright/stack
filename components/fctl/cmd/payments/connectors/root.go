package connectors

import (
	"github.com/formancehq/fctl/cmd/payments/connectors/install"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewConnectorsCommand() *cobra.Command {
	return fctl.NewCommand("connectors",
		fctl.WithAliases("c", "co", "con"),
		fctl.WithShortDescription("Connectors management"),
		fctl.WithChildCommands(
			NewGetConfigCommand(),
			NewUninstallCommand(),
			NewListCommand(),
			install.NewInstallCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
