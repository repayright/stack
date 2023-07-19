package fctl

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
)

type Renderable interface {
	Render() error
}

type Controller[T any] interface {
	GetStore() T

	GetConfig() *ControllerConfig

	Run() (Renderable, error)
}
type ExportedData struct {
	Data interface{} `json:"data"`
}

type FValue struct {
	name string
}

func (f *FValue) String() string {
	return f.name
}

func (f *FValue) Set(string) error {
	f.name = f.name
	return nil
}

var (
	stackFlagV        *FValue   = &FValue{name: ""}
	organizationFlagV *FValue   = &FValue{name: ""}
	ledgerFlagV       *FValue   = &FValue{name: ""}
	Stack             flag.Flag = flag.Flag{
		Name:     "stack",
		Usage:    "stack name",
		DefValue: "",
		Value:    stackFlagV,
	}
	Organization flag.Flag = flag.Flag{
		Name:     "organization",
		Usage:    "organization name",
		DefValue: "",
		Value:    organizationFlagV,
	}
	Ledger flag.Flag = flag.Flag{
		Name:     "ledger",
		Usage:    "Specific ledger",
		DefValue: "default",
		Value:    ledgerFlagV,
	}
)

type ControllerConfig struct {
	context          context.Context
	use              string
	description      string
	shortDescription string
	aliases          []string
	out              io.Writer
	flags            *flag.FlagSet
	pflags           *flag.FlagSet
	scope            *flag.FlagSet
	args             []string
}

var GlobalFlags *flag.FlagSet = func() *flag.FlagSet {
	flags := flag.NewFlagSet("global", flag.ContinueOnError)

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flags.Bool(InsecureTlsFlag, false, "insecure TLS")
	flags.Bool(TelemetryFlag, false, "enable telemetry")
	flags.Bool(DebugFlag, false, "debug mode")
	flags.String(ProfileFlag, "", "config profile to use")
	flags.String(FileFlag, fmt.Sprintf("%s/.formance/fctl.config", homedir), "config file to use")
	flags.String(outputFlag, "plain", "output format (plain, json)")

	return flags
}()

func generateScopesEnum(s ...flag.Flag) *flag.FlagSet {
	fs := flag.NewFlagSet("scopes", flag.ExitOnError)

	if len(s) == 0 {
		return fs
	}
	for _, f := range s {
		fmt.Println(f.Value)
		fs.Var(f.Value, f.Name, f.Usage)
	}
	return fs
}

func NewControllerConfig(use string, description string, shortDescription string, aliases []string, out io.Writer, flags *flag.FlagSet, s ...flag.Flag) *ControllerConfig {

	return &ControllerConfig{
		use:              use,
		description:      description,
		shortDescription: shortDescription,
		aliases:          aliases,
		out:              out,
		flags:            flags,
		scope:            generateScopesEnum(s...),
		pflags:           GlobalFlags,
	}

}

func (c *ControllerConfig) GetScopes() *flag.FlagSet {
	return c.scope
}

func (c *ControllerConfig) GetUse() string {
	return c.use
}

func (c *ControllerConfig) GetDescription() string {
	return c.description
}

func (c *ControllerConfig) GetShortDescription() string {
	return c.shortDescription
}

func (c *ControllerConfig) SetShortDescription(shortDescription string) {
	c.shortDescription = shortDescription
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
func (c *ControllerConfig) SetOut(out io.Writer) {
	c.out = out
}

func (c *ControllerConfig) GetArgs() []string {
	return c.args
}

func (c *ControllerConfig) SetArgs(args []string) {
	c.args = args
}

func (c *ControllerConfig) GetFlags() *flag.FlagSet {
	return c.flags
}

// GetAllFLags Return the pflags & flags merged together in a new FlagSet
// This is done to avoid mutating the original flag.FlagSet
// which is used by the controller to parse the flags
// and the pflags are used by the controller to parse the persistent one
func (c *ControllerConfig) GetAllFLags() *flag.FlagSet {

	// Create a new FlagSet
	flags := flag.NewFlagSet(c.use, flag.ExitOnError)

	// Regroup pflag in 1 flagset
	if c.pflags != nil {
		c.pflags.VisitAll(func(f *flag.Flag) {
			flags.Var(f.Value, f.Name, f.Usage)
		})
	}

	// Regroup flags in 1 flagset
	if c.flags != nil {
		c.flags.VisitAll(func(f *flag.Flag) {
			flags.Var(f.Value, f.Name, f.Usage)
		})
	}

	return flags
}

func (c *ControllerConfig) GetPFlags() *flag.FlagSet {
	return c.pflags
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
