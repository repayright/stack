package modelutils

import (
	"context"
	"flag"
)

// type Renderable interface {
// 	Render(cmd *cobra.Command, args []string) (ui.Model, error)
// }
// type Controller[T any] interface {
// 	GetStore() T
// 	Run(cmd *cobra.Command, args []string) (Renderable, error)
// }

type Renderable interface {
	Render() (Model, error)
}

type ControllerContext struct {
	flags   *flag.FlagSet
	args    []string
	context context.Context
}

type Controller[T any] interface {
	GetStore() T
	GetKeyMapAction() *KeyMapHandler[T]

	/**
	 * The following methods are used to render the UI
	 * Implementing a cycle of Init, Run, Render
	 *
	 **/
	// Use to populate contextual config
	Init()

	// Use to render to a defined UI model
	Render() (Model, error)

	// Use to populate the store
	// It depends on the controller and the cobra configuration passed to the controller
	// Who is responsible for populating the store and act depending on flags and args
	Run() error
}
type ExportedData struct {
	Data interface{} `json:"data"`
}
