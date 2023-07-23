package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/invitations"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("organizations",
		fctl.WithAliases("org", "o"),
		fctl.WithShortDescription("Organizations management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
			invitations.NewCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization, config.Stack),
	)
}
