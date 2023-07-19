package auth

import (
	"github.com/formancehq/fctl/cmd/auth/clients"
	"github.com/formancehq/fctl/cmd/auth/users"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("auth",
		fctl.WithShortDescription("Auth server management"),
		fctl.WithChildCommands(
			clients.NewCommand(),
			users.NewCommand(),
		),
		fctl.WithScopesFlags(fctl.Organization, fctl.Stack),
	)
}
