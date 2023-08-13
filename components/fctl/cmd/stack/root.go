package stack

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewNodeController() *config.Node {
	return config.NewConfigNode(*config.NewControllerConfig(
		"stack",
		"Manage your stack",
		"stacks",
		[]string{"stack", "stacks", "st"},
		nil,
	),
		NewListController(NewListControllerConfig()),
	)
}

func NewCommand() *cobra.Command {
	return fctl.NewCommand("stack",
		fctl.WithShortDescription("Manage your stack"),
		fctl.WithAliases("stack", "stacks", "st"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewListCommand(),
			NewDeleteCommand(),
			NewShowCommand(),
			NewRestoreStackCommand(),
		),
		fctl.WithCommandScopesFlags(config.Organization),
	)
}
