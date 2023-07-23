package config

import (
	"flag"
	"sync"
)

var scopeFlags *flag.FlagSet
var lock = &sync.Mutex{}

var (
	stackFlagV = &fValue[string]{
		value: "",
	}
	organizationFlagV = &fValue[string]{
		value: "",
	}
	ledgerFlagV = &fValue[string]{
		value: "default",
	}
)

var (
	Stack        = getScopeFlags(StackFlag)
	Ledger       = getScopeFlags("ledger")
	Organization = getScopeFlags(OrganizationFlag)
)

func getScopesFlagsInstance() *flag.FlagSet {
	if scopeFlags == nil {
		lock.Lock()
		defer lock.Unlock()
		if scopeFlags == nil {
			scopeFlags = flag.NewFlagSet("scopes", flag.ContinueOnError)
			scopeFlags.StringVar(stackFlagV.Get(), "stack", "", "Specific stack id (not required if only one stack is present)")
			scopeFlags.StringVar(organizationFlagV.Get(), "organization", "", "Selected organization (not required if only one organization is present)")
			scopeFlags.StringVar(ledgerFlagV.Get(), "ledger", "default", "Specific ledger name")
		}
	}

	return scopeFlags
}

func getScopeFlags(name string) *flag.Flag {
	return getScopesFlagsInstance().Lookup(name)
}
