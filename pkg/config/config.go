package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
	"github.com/LeMyst/kube-no-trouble/pkg/printer"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
)

type Config struct {
	AdditionalKinds       []string
	AdditionalAnnotations []string
	Cluster               bool
	Context               string
	ExitError             bool
	Filenames             []string
	Helm3                 bool
	Kubeconfig            string
	LogLevel              ZeroLogLevel
	Output                string
	OutputFile            string
	TargetVersion         *judge.Version
	KubentVersion         bool
}

func NewFromFlags() (*Config, error) {
	config := Config{
		LogLevel:      ZeroLogLevel(zerolog.InfoLevel),
		TargetVersion: &judge.Version{},
	}

	flag.StringSliceVarP(&config.AdditionalKinds, "additional-kind", "a", []string{}, "additional kinds of resources to report in Kind.version.group.com format")
	flag.StringSliceVarP(&config.AdditionalAnnotations, "additional-annotation", "A", []string{}, "additional annotations that should be checked to determine the last applied config")
	flag.BoolVarP(&config.Cluster, "cluster", "c", true, "enable Cluster collector")
	flag.StringVarP(&config.Context, "context", "x", "", "kubeconfig context")
	flag.BoolVarP(&config.ExitError, "exit-error", "e", false, "exit with non-zero code when issues are found")
	flag.BoolVarP(&config.KubentVersion, "version", "v", false, "prints the version of kubent and exits")
	flag.BoolVar(&config.Helm3, "helm3", true, "enable Helm v3 collector")
	flag.StringSliceVarP(&config.Filenames, "filename", "f", []string{}, "manifests to check, use - for stdin")
	flag.StringVarP(&config.Kubeconfig, "kubeconfig", "k", os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "path to the kubeconfig file")
	flag.StringVarP(&config.Output, "output", "o", "text", "output format - [text|json|csv]")
	flag.StringVarP(&config.OutputFile, "output-file", "O", "-", "output file, use - for stdout")
	flag.VarP(&config.LogLevel, "log-level", "l", "set log level (trace, debug, info, warn, error, fatal, panic, disabled)")
	flag.VarP(config.TargetVersion, "target-version", "t", "target K8s version in SemVer format (autodetected by default)")

	flag.Parse()

	if _, err := printer.ParsePrinter(config.Output); err != nil {
		return nil, fmt.Errorf("failed to validate argument output: %w", err)
	}

	if err := validateOutputFile(config.OutputFile); err != nil {
		return nil, fmt.Errorf("failed to validate argument output-file: %w", err)
	}

	if err := validateAdditionalResources(config.AdditionalKinds); err != nil {
		return nil, fmt.Errorf("failed to validate arguments: %w", err)
	}

	// This is a little ugly, but I think preferred to implementing
	// unset semantics & logic compared to using nil
	// and should be solvable by using new https://pkg.go.dev/flag#Func
	if config.TargetVersion.Version == nil {
		config.TargetVersion = nil
	}

	return &config, nil
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

// validateOutputFile checks if output file name is valid and if the
// destination directory exists
func validateOutputFile(outputFileName string) error {
	if outputFileName == "" {
		return fmt.Errorf("output file name can't be empty (use - for stdout)")
	}

	if outputFileName != "-" {
		dir := filepath.Dir(outputFileName)
		if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("output directory %s does not exist", dir)
		}
	}

	return nil
}
