package prompt

import (
	"reflect"

	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/fctl/pkg/helpers"
)

type Commands struct {
	node          *config.Node
	commands      []string
	descmap       map[string]string
	controllerMap map[string]config.Controller
}

func NewCommands(node *config.Node) *Commands {
	Log := helpers.NewLogger("prompt")
	Log.Log("NewCommands")
	var commands []string
	var description map[string]string = make(map[string]string)
	var controllerMap map[string]config.Controller = make(map[string]config.Controller)

	childs := node.GetChilds()

	if childs == nil {
		return &Commands{
			commands: commands,
			descmap:  description,
		}
	}

	typ := reflect.TypeOf(childs).String()
	Log.Log(typ)

	switch c := childs.(type) {
	case []*config.Node:
		for _, child := range c {
			childsCommand := NewCommands(child)
			commands = append(commands, childsCommand.commands...)
			for k, v := range childsCommand.descmap {
				description[k] = v
			}

			for k, v := range childsCommand.controllerMap {
				controllerMap[k] = v
			}
		}
	case *config.ConfigNode:
		conf := c.Config
		for _, child := range c.Controllers {
			node := config.NewControllerNode(child)
			newCommands := NewCommands(node)

			for _, command := range newCommands.commands {
				key := conf.GetUse() + " " + command
				commands = append(commands, key)
				description[key] = newCommands.descmap[command]
				controllerMap[key] = child
			}

		}
	case []config.Controller:
		for _, child := range c {
			config := child.GetConfig()
			commands = append(commands, config.GetUse())
			description[config.GetUse()] = config.GetShortDescription()
			controllerMap[config.GetUse()] = child
		}
	}

	Log.Log(commands...)
	return &Commands{
		node:          node,
		commands:      commands,
		descmap:       description,
		controllerMap: controllerMap,
	}
}
