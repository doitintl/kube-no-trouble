package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type ZeroLogLevel int8

func (l ZeroLogLevel) String() string {
	return zerolog.Level(l).String()
}

func (l *ZeroLogLevel) Set(s string) error {
	level, err := zerolog.ParseLevel(s)
	if err != nil {
		return err
	}
	*l = ZeroLogLevel(level)
	return nil
}

func (l ZeroLogLevel) Type() string {
	return "string"
}

type Config struct {
	AdditionalKinds []string
	Cluster         bool
	ExitError       bool
	Filenames       []string
	Helm2           bool
	Helm3           bool
	Kubeconfig      string
	LogLevel        ZeroLogLevel
	Output          string
}

func NewFromFlags() (*Config, error) {
	config := Config{LogLevel: ZeroLogLevel(zerolog.InfoLevel)}

	home := homedir.HomeDir()
	flag.StringSliceVarP(&config.AdditionalKinds, "additional-kind", "a", []string{}, "additional kinds of resources to report in Kind.version.group.com format")
	flag.BoolVarP(&config.Cluster, "cluster", "c", true, "enable Cluster collector")
	flag.BoolVarP(&config.ExitError, "exit-error", "e", false, "exit with non-zero code when issues are found")
	flag.BoolVar(&config.Helm2, "helm2", true, "enable Helm v2 collector")
	flag.BoolVar(&config.Helm3, "helm3", true, "enable Helm v3 collector")
	flag.StringSliceVarP(&config.Filenames, "filename", "f", []string{}, "manifests to check, use - for stdin")
	flag.StringVarP(&config.Kubeconfig, "kubeconfig", "k", envOrString("KUBECONFIG", filepath.Join(home, ".kube", "config")), "path to the kubeconfig file")
	flag.StringVarP(&config.Output, "output", "o", "text", "output format - [text|json]")
	flag.VarP(&config.LogLevel, "log-level", "l", "set log level (trace, debug, info, warn, error, fatal, panic, disabled)")

	flag.Parse()

	if err := validateAdditionalResources(config.AdditionalKinds); err != nil {
		return nil, fmt.Errorf("failed to validate arguments: %w", err)
	}

	return &config, nil
}

func envOrString(env string, def string) string {
	val, ok := os.LookupEnv(env)
	if ok {
		return val
	}
	return def
}

// validateAdditionalResources check that all resources are provided in full form
// resource.version.group.com. E.g. managedcertificate.v1beta1.networking.gke.io
func validateAdditionalResources(resources []string) error {
	for _, r := range resources {
		parts := strings.Split(r, ".")
		if len(parts) < 4 {
			return fmt.Errorf("failed to parse additional Kind, full form Kind.version.group.com is expected, instead got: %s", r)
		}

		if !unicode.IsUpper(rune(parts[0][0])) {
			return fmt.Errorf("failed to parse additional Kind, Kind is expected to be capitalized by convention, instead got: %s", parts[0])
		}
	}
	return nil
}
