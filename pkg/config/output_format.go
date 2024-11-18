package config

import "fmt"

const (
	JSON = "json"
	TEXT = "text"
	CSV  = "csv"
)

// implements pflag.Value to enforce strict string matching for the output format during flag parsing
type OutputFormat string

func (o *OutputFormat) String() string {
	return string(*o)
}

func (o *OutputFormat) Set(value string) error {
	switch value {
	case JSON, TEXT, CSV:
		*o = OutputFormat(value)
		return nil
	default:
		return fmt.Errorf("invalid output format: %s (must be one of: json, text, csv)", value)
	}
}

func (o *OutputFormat) Type() string {
	return "OutputFormat"
}
