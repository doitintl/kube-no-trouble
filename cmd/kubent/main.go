package main

import (
	"fmt"
	"os"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
	"github.com/doitintl/kube-no-trouble/pkg/config"
	"github.com/doitintl/kube-no-trouble/pkg/judge"
	"github.com/doitintl/kube-no-trouble/pkg/printer"
	"github.com/doitintl/kube-no-trouble/pkg/rules"

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
	EXIT_CODE_SUCCESS      = 0
	EXIT_CODE_FAIL_GENERIC = 1
	EXIT_CODE_FOUND_ISSUES = 200
)

func getCollectors(collectors []collector.Collector) []map[string]interface{} {
	var inputs []map[string]interface{}
	for _, c := range collectors {
		rs, err := c.Get()
		if err != nil {
			log.Error().Err(err).Str("name", c.Name()).Msg("Failed to retrieve data from collector")
		} else {
			inputs = append(inputs, rs...)
			log.Info().Str("name", c.Name()).Msgf("Retrieved %d resources from collector", len(rs))
		}
	}
	return inputs
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
		collector, err := collector.NewClusterCollector(&collector.ClusterOpts{Kubeconfig: config.Kubeconfig}, config.AdditionalKinds)
		collectors = storeCollector(collector, err, collectors)
	}

	if config.Helm2 {
		collector, err := collector.NewHelmV2Collector(&collector.HelmV2Opts{Kubeconfig: config.Kubeconfig})
		collectors = storeCollector(collector, err, collectors)
	}

	if config.Helm3 {
		collector, err := collector.NewHelmV3Collector(&collector.HelmV3Opts{Kubeconfig: config.Kubeconfig})
		collectors = storeCollector(collector, err, collectors)
	}

	if len(config.Filenames) > 0 {
		collector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: config.Filenames})
		collectors = storeCollector(collector, err, collectors)
	}
	return collectors
}

func getServerVersion(cv *config.Version, collectors []collector.Collector) error {
	if cv.Version == nil {
		for _, c := range collectors {
			if versionCol, ok := c.(collector.VersionCollector); ok {
				goversion, err := versionCol.GetServerVersion()
				if err != nil {
					return fmt.Errorf("failed to detect k8s version: %w", err)
				}

				cv.SetFromVersion(goversion)
				return nil
			}
		}
	}
	return nil
}

func main() {
	exitCode := EXIT_CODE_FAIL_GENERIC

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config, err := config.NewFromFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config flags")
	}

	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))

	log.Info().Msg(">>> Kube No Trouble `kubent` <<<")
	log.Info().Msgf("version %s (git sha %s)", version, gitSha)

	log.Info().Msg("Initializing collectors and retrieving data")
	initCollectors := initCollectors(config)

	getServerVersion(&config.TargetVersion, initCollectors)
	if config.TargetVersion.Version != nil {
		log.Info().Msgf("Target K8s version is %s", config.TargetVersion.String())
	}

	collectors := getCollectors(initCollectors)

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

	printer, err := printer.NewPrinter(config.Output)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create printer")
	}

	err = printer.Print(results)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to print results")
	}

	if config.ExitError && len(results) > 0 {
		exitCode = EXIT_CODE_FOUND_ISSUES
	} else {
		exitCode = EXIT_CODE_SUCCESS
	}
	os.Exit(exitCode)
}
