package cloud

import (
	"github.com/formancehq/fctl/cmd/cloud/billing"
	"github.com/formancehq/fctl/cmd/cloud/me"
	"github.com/formancehq/fctl/cmd/cloud/organizations"
	"github.com/formancehq/fctl/cmd/cloud/regions"
	"github.com/formancehq/fctl/cmd/cloud/users"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("cloud",
		fctl.WithAliases("c"),
		fctl.WithShortDescription("Cloud management"),
		fctl.WithChildCommands(
			organizations.NewCommand(),
			me.NewCommand(),
			users.NewCommand(),
			regions.NewCommand(),
			billing.NewCommand(),
			NewGeneratePersonalTokenCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
