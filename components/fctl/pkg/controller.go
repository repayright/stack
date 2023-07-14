package fctl

import (
	"context"
	"flag"
	"io"
	"os"
)

type Renderable interface {
	Render() error
}
type Controller[T any] interface {
	GetStore() T

	GetConfig() ControllerConfig

	Run() (Renderable, error)
}
type ExportedData struct {
	Data interface{} `json:"data"`
}

type ControllerConfig struct {
	context          context.Context
	use              string
	description      string
	shortDescription *string
	aliases          []string
	out              io.Writer
	flags            *flag.FlagSet
	args             []string
}

func NewControllerConfig(use string, description string, aliases []string, out io.Writer, flags *flag.FlagSet) *ControllerConfig {
	WithGlobalFlags(flags)
	return &ControllerConfig{
		use:         use,
		description: description,
		aliases:     aliases,
		out:         out,
		flags:       flags,
	}

}
func (c *ControllerConfig) GetUse() string {
	return c.use
}

func (c *ControllerConfig) GetDescription() string {
	return c.description
}

func (c *ControllerConfig) GetShortDescription() *string {
	return c.shortDescription
}

func (c *ControllerConfig) SetShortDescription(shortDescription string) {
	c.shortDescription = &shortDescription
}

func (c *ControllerConfig) GetAliases() []string {
	return c.aliases
}

func (c *ControllerConfig) GetOut() io.Writer {
	if c.out == nil {
		return os.Stdout
	}

	return c.out
}

func (c *ControllerConfig) GetArgs() []string {
	return c.args
}

func (c *ControllerConfig) SetArgs(args []string) {
	c.args = append([]string{}, args...)
}

func (c *ControllerConfig) GetFlags() *flag.FlagSet {
	return c.flags
}

func (c *ControllerConfig) GetContext() context.Context {
	if c.context == nil {
		return context.TODO()
	}

	return c.context
}

func (c *ControllerConfig) SetContext(ctx context.Context) {
	c.context = ctx
}
