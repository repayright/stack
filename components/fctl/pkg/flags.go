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

func GetBool(flags *flag.FlagSet, flagName string) bool {
	f := flags.Lookup(flagName)
	if f == nil {
		return false
	}

	value := f.Value.String()
	if value == "" {
		return false
	}
	v, err := strconv.ParseBool(value)

	if err != nil {
		fromEnv := strings.ToLower(os.Getenv(strcase.ToScreamingSnake(flagName)))
		return fromEnv == "true" || fromEnv == "1"
	}
	return v
}

func GetString(flagSet *flag.FlagSet, flagName string) string {
	f := flagSet.Lookup(flagName)
	if f == nil {
		return ""
	}

	v := f.Value.String()
	if v == "" {
		return os.Getenv(strcase.ToScreamingSnake(flagName))
	}
	return v
}

func GetStringSlice(flagSet *flag.FlagSet, flagName string) []string {

	f := flagSet.Lookup(flagName)
	if f == nil {
		return []string{}
	}

	value := f.Value.String()
	v := strings.Split(value, " ")
	if value == "" {
		return []string{}
	}

	if len(v) == 0 {
		envVar := os.Getenv(strcase.ToScreamingSnake(flagName))
		if envVar == "" {
			return []string{}
		}
		return strings.Split(envVar, " ")
	}
	return v
}

func GetInt(flagSet *flag.FlagSet, flagName string) int {

	f := flagSet.Lookup(flagName)
	if f == nil {
		return 0
	}

	value := f.Value.String()
	if value == "" {
		return 0
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		v := os.Getenv(strcase.ToScreamingSnake(flagName))
		if v != "" {
			v, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return 0
			}
			return int(v)
		}
		return 0
	}
	return int(v)
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
