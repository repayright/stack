package fctl

import "flag"

type fValue struct {
	name string
}

func (f *fValue) String() string {
	return f.name
}

func (f *fValue) Set(s string) error {
	f.name = s
	return nil
}

func (f *fValue) Get() *string {
	return &f.name
}

var (
	stackFlagV        *fValue   = &fValue{name: ""}
	organizationFlagV *fValue   = &fValue{name: ""}
	ledgerFlagV       *fValue   = &fValue{name: ""}
	Stack             flag.Flag = flag.Flag{
		Name:     "stack",
		Usage:    "stack id",
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
