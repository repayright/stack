package prompt

import (
	"github.com/spf13/cobra"
)

type Commands struct {
	commands []string
	descmap  map[string]string
}

func NewCommands(cmd *cobra.Command) *Commands {
	var commands []string
	var description map[string]string = make(map[string]string)
	for _, c := range cmd.Commands() {
		commands = append(commands, c.Use)
		description[c.Use] = c.Short
		if c.HasSubCommands() {
			subCommands := NewCommands(c)
			// Prefix sub tree with parent command
			for i, sub := range subCommands.commands {
				subCommands.commands[i] = c.Use + " " + sub
				description[subCommands.commands[i]] = subCommands.descmap[sub]
			}

			commands = append(commands, subCommands.commands...)
		}
	}

	return &Commands{
		commands: commands,
		descmap:  description,
	}
}
