package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LeMyst/kube-no-trouble/pkg/collector"
	"github.com/LeMyst/kube-no-trouble/pkg/config"
	"github.com/LeMyst/kube-no-trouble/pkg/judge"
	"github.com/LeMyst/kube-no-trouble/pkg/printer"
	"github.com/LeMyst/kube-no-trouble/pkg/rules"

	"k8s.io/klog/v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var (
	version string = "dev"
	gitSha  string = "dev"
)

const (
	EXIT_CODE_SUCCESS             = 0
	EXIT_CODE_FAIL_GENERIC        = 1
	EXIT_CODE_FAIL_GET_COLLECTORS = 100
	EXIT_CODE_FOUND_ISSUES        = 200
)

func generateUserAgent() string {
	return fmt.Sprintf("kubent (%s/%s)", version, gitSha)
}

func getCollectors(collectors []collector.Collector) ([]map[string]interface{}, error) {
	var inputs []map[string]interface{}
	for _, c := range collectors {
		rs, err := c.Get()
		if err != nil {
			log.Error().Err(err).Str("name", c.Name()).Msg("Failed to retrieve data from collector")
			return nil, err
		}
		inputs = append(inputs, rs...)
		log.Info().Str("name", c.Name()).Msgf("Retrieved %d resources from collector", len(rs))
	}
	return inputs, nil
}

func storeCollector(collector collector.Collector, err error, collectors []collector.Collector) []collector.Collector {
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to initialize collector: %+v", collector))
	} else {
		collectors = append(collectors, collector)
	}
	return collectors
}

func initCollectors(config *config.Config) []collector.Collector {
	collectors := []collector.Collector{}
	if config.Cluster {
		collector, err := collector.NewClusterCollector(&collector.ClusterOpts{Kubeconfig: config.Kubeconfig, KubeContext: config.Context},
			config.AdditionalKinds, config.AdditionalAnnotations, generateUserAgent())
		collectors = storeCollector(collector, err, collectors)
	}

	if config.Helm3 {
		collector, err := collector.NewHelmV3Collector(&collector.HelmV3Opts{Kubeconfig: config.Kubeconfig, KubeContext: config.Context}, generateUserAgent())
		collectors = storeCollector(collector, err, collectors)
	}

	if len(config.Filenames) > 0 {
		collector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: config.Filenames})
		collectors = storeCollector(collector, err, collectors)
	}
	return collectors
}

func getServerVersion(cv *judge.Version, collectors []collector.Collector) (*judge.Version, error) {
	if cv == nil {
		for _, c := range collectors {
			if versionCol, ok := c.(collector.VersionCollector); ok {
				version, err := versionCol.GetServerVersion()
				if err != nil {
					return nil, fmt.Errorf("failed to detect k8s version: %w", err)
				}
				return version, nil
			}
		}
	}
	return cv, nil
}

func main() {
	exitCode := EXIT_CODE_FAIL_GENERIC

	configureGlobalLogging()

	config, err := config.NewFromFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config flags")
	}

	if config.KubentVersion {
		log.Info().Msgf("version %s (git sha %s)", version, gitSha)
		os.Exit(EXIT_CODE_SUCCESS)
	}

	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))

	log.Info().Msg(">>> Kube No Trouble `kubent` <<<")
	log.Info().Msgf("version %s (git sha %s)", version, gitSha)

	log.Info().Msg("Initializing collectors and retrieving data")
	initCollectors := initCollectors(config)

	config.TargetVersion, _ = getServerVersion(config.TargetVersion, initCollectors)
	if config.TargetVersion != nil {
		log.Info().Msgf("Target K8s version is %s", config.TargetVersion.String())
	}

	collectors, err := getCollectors(initCollectors)
	if err != nil {
		os.Exit(EXIT_CODE_FAIL_GET_COLLECTORS)
	}

	// this could probably use some error checking in future, but
	// schema.ParseKindArg does not return any error
	var additionalKinds []schema.GroupVersionKind
	for _, ar := range config.AdditionalKinds {
		gvr, _ := schema.ParseKindArg(ar)
		additionalKinds = append(additionalKinds, *gvr)
	}

	loadedRules, err := rules.FetchRegoRules(additionalKinds)
	if err != nil {
		log.Fatal().Err(err).Str("name", "Rules").Msg("Failed to load rules")
	}

	judge, err := judge.NewRegoJudge(&judge.RegoOpts{}, loadedRules)
	if err != nil {
		log.Fatal().Err(err).Str("name", "Rego").Msg("Failed to initialize decision engine")
	}

	results, err := judge.Eval(collectors)
	if err != nil {
		log.Fatal().Err(err).Str("name", "Rego").Msg("Failed to evaluate input")
	}

	results, err = printer.FilterNonRelevantResults(results, config.TargetVersion)
	if err != nil {
		log.Fatal().Err(err).Str("name", "Rego").Msg("Failed to filter results")
	}

	err = outputResults(results, config.Output, config.OutputFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to output results")
	}

	if config.ExitError && len(results) > 0 {
		exitCode = EXIT_CODE_FOUND_ISSUES
	} else {
		exitCode = EXIT_CODE_SUCCESS
	}
	os.Exit(exitCode)
}

func configureGlobalLogging() {
	// disable any logging from K8S go-client
	// unfortunately restConfig.WarningHandler does not handle auth plugins
	klog.SetOutput(io.Discard)
	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	flags.Set("logtostderr", "false")

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func outputResults(results []judge.Result, outputType string, outputFile string) error {
	printer, err := printer.NewPrinter(outputType, outputFile)
	if err != nil {
		return fmt.Errorf("failed to create printer: %v", err)
	}
	defer printer.Close()

	err = printer.Print(results)
	if err != nil {
		return fmt.Errorf("failed to print results: %v", err)
	}

	return nil
}
