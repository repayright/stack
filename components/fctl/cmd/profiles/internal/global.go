package internal

import (
	"flag"

	"github.com/formancehq/fctl/pkg/config"
	"github.com/spf13/cobra"
)

func ProfileNamesAutoCompletion(flags *flag.FlagSet, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ret, err := config.ListProfiles(flags, toComplete)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}

	return ret, cobra.ShellCompDirectiveDefault
}

func ProfileCobraAutoCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	flags := config.ConvertPFlagSetToFlagSet(cmd.Flags())
	return ProfileNamesAutoCompletion(flags, args, toComplete)
}
