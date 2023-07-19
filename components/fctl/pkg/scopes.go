package fctl

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type fValue[T any] struct {
	value T
}

func (f *fValue[T]) String() T {
	return f.value
}

func (f *fValue[T]) Set(s T) error {
	f.value = s
	return nil
}

func (f *fValue[T]) Get() *T {
	return &f.value
}

var (
	stackFlagV        *fValue[string] = &fValue[string]{value: ""}
	organizationFlagV *fValue[string] = &fValue[string]{value: ""}
	ledgerFlagV       *fValue[string] = &fValue[string]{value: ""}
	insecureTlsV      *fValue[bool]   = &fValue[bool]{value: false}
	telemetryFlagV    *fValue[bool]   = &fValue[bool]{value: false}
	debugFlagV        *fValue[bool]   = &fValue[bool]{value: false}
	profileFlagV      *fValue[string] = &fValue[string]{value: ""}
	configFlagV       *fValue[string] = &fValue[string]{value: fmt.Sprintf("%s/.formance/fctl.config", getHomeDir())}
	outputFlagV       *fValue[string] = &fValue[string]{value: "plain"}
	Stack             flag.Flag       = flag.Flag{
		Name:     "stack",
		Usage:    "Specific stack (not required if only one stack is present)",
		DefValue: "",
		Value:    stackFlagV,
	}
	Organization flag.Flag = flag.Flag{
		Name:     "organization",
		Usage:    "Selected organization (not required if only one organization is present)",
		DefValue: "",
		Value:    organizationFlagV,
	}
	Ledger flag.Flag = flag.Flag{
		Name:     "ledger",
		Usage:    "Specific ledger name",
		DefValue: "default",
		Value:    ledgerFlagV,
	}
)

var GlobalFlags *flag.FlagSet = func() *flag.FlagSet {
	flags := flag.NewFlagSet("global", flag.ContinueOnError)
	flags.BoolVar(insecureTlsV.Get(), InsecureTlsFlag, false, "insecure TLS")
	flags.BoolVar(telemetryFlagV.Get(), TelemetryFlag, false, "enable telemetry")
	flags.BoolVar(debugFlagV.Get(), DebugFlag, false, "debug mode")
	flags.StringVar(profileFlagV.Get(), ProfileFlag, "", "config profile to use")
	flags.StringVar(configFlagV.Get(), ConfigFlag, fmt.Sprintf("%s/.formance/fctl.config", getHomeDir()), "config file to use")
	flags.StringVar(outputFlagV.Get(), outputFlag, "plain", "output format (plain, json)")

	return flags
}()

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(errors.New("unable to get home directory"))
	}
	return homedir
}
