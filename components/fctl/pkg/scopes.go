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
	stackFlagV        = &fValue[string]{value: ""}
	organizationFlagV = &fValue[string]{value: ""}
	ledgerFlagV       = &fValue[string]{value: ""}
	insecureTlsV      = &fValue[bool]{value: false}
	telemetryFlagV    = &fValue[bool]{value: false}
	debugFlagV        = &fValue[bool]{value: false}
	profileFlagV      = &fValue[string]{value: ""}
	configFlagV       = &fValue[string]{value: fmt.Sprintf("%s/.formance/fctl.config", getHomeDir())}
	outputFlagV       = &fValue[string]{value: "plain"}
	scopeFlags        = func() *flag.FlagSet {
		flags := flag.NewFlagSet("scopes", flag.ContinueOnError)
		flags.StringVar(stackFlagV.Get(), stackFlag, "", "Specific stack id (not required if only one stack is present)")
		flags.StringVar(organizationFlagV.Get(), organizationFlag, "", "Selected organization (not required if only one organization is present)")
		flags.StringVar(ledgerFlagV.Get(), "ledger", "", "Specific ledger name")

		return flags
	}()
	Stack        = getScopeFlags(stackFlag)
	Ledger       = getScopeFlags("ledger")
	Organization = getScopeFlags(organizationFlag)
)

var GlobalFlags = func() *flag.FlagSet {
	flags := flag.NewFlagSet("global", flag.ContinueOnError)
	flags.BoolVar(insecureTlsV.Get(), InsecureTlsFlag, false, "insecure TLS")
	flags.BoolVar(telemetryFlagV.Get(), TelemetryFlag, false, "enable telemetry")
	flags.BoolVar(debugFlagV.Get(), DebugFlag, false, "debug mode")
	flags.StringVar(profileFlagV.Get(), ProfileFlag, "", "config profile to use")
	flags.StringVar(configFlagV.Get(), ConfigFlag, fmt.Sprintf("%s/.formance/fctl.config", getHomeDir()), "config file to use")
	flags.StringVar(outputFlagV.Get(), outputFlag, "plain", "output format (plain, json)")

	return flags
}()

func getScopeFlags(name string) *flag.Flag {
	return scopeFlags.Lookup(name)
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(errors.New("unable to get home directory"))
	}
	return homedir
}
