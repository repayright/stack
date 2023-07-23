package install

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	return fctl.NewCommand("install",
		fctl.WithAliases("i"),
		fctl.WithShortDescription("Install a connector"),
		fctl.WithChildCommands(
			NewStripeCommand(),
			NewBankingCircleCommand(),
			NewCurrencyCloudCommand(),
			NewModulrCommand(),
			NewWiseCommand(),
			NewMangoPayCommand(),
			NewMoneycorpCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
