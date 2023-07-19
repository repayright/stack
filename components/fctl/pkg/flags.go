package fctl

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
)

const (
	MembershipURIFlag string = "membership-uri"
	FileFlag          string = "config"
	ProfileFlag       string = "profile"
	OutputFlag        string = "output"
	DebugFlag         string = "debug"
	InsecureTlsFlag   string = "insecure-tls"
	TelemetryFlag     string = "telemetry"
	MetadataFlag      string = "metadata"
)

func WithStackPersistentFlag(flag *flag.FlagSet, name, defaultValue, help string) *flag.FlagSet {
	flag.String(stackFlag, "", "Specific stack (not required if only one stack is present)")
	return flag
}
func GetBool(flags *flag.FlagSet, flagName string) bool {
	f := flags.Lookup(flagName)
	if f == nil {
		return false
	}

	fromEnv := strings.ToLower(os.Getenv(strcase.ToScreamingSnake(flagName)))
	if fromEnv != "" {
		return fromEnv == "true" || fromEnv == "1"
	}

	value := f.Value.String()
	if value == "" {
		return false
	}
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return v
}

func GetString(flagSet *flag.FlagSet, flagName string) string {
	f := flagSet.Lookup(flagName)
	if f == nil {
		return ""
	}

	envVar := os.Getenv(strcase.ToScreamingSnake(flagName))
	if envVar != "" {
		return envVar
	}

	return f.Value.String()
}

func GetStringSlice(flagSet *flag.FlagSet, flagName string) []string {

	f := flagSet.Lookup(flagName)

	if f == nil {
		return []string{}
	}

	envVar := os.Getenv(strcase.ToScreamingSnake(flagName))

	if len(envVar) > 0 {
		return strings.Split(envVar, " ")
	}

	value := f.Value.String()
	if len(value) > 0 {
		return strings.Split(value, " ")
	}

	return []string{}
}

func GetInt(flagSet *flag.FlagSet, flagName string) int {

	f := flagSet.Lookup(flagName)
	if f == nil {
		return 0
	}

	v := os.Getenv(strcase.ToScreamingSnake(flagName))
	if v != "" {
		v, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0
		}
		return int(v)
	}

	value := f.Value.String()
	if value == "" {
		return 0
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return 0
	}
	return int(i)
}

func WithConfirmFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("confirm", false, "Confirm the action")
}

func WithMetadataFlag(flag *flag.FlagSet) *flag.FlagSet {
	flag.String(MetadataFlag, "", "Metadata to use")
	return flag
}
func ConvertPFlagSetToFlagSet(pFlagSet *pflag.FlagSet) *flag.FlagSet {

	flagSet := flag.NewFlagSet("fctl", flag.ExitOnError)

	pFlagSet.VisitAll(func(f *pflag.Flag) {
		flagSet.Var(f.Value, f.Name, f.Usage)
	})

	return flagSet
}
