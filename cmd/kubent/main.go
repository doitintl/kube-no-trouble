package main

import (
	"os"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
	"github.com/doitintl/kube-no-trouble/pkg/config"
	"github.com/doitintl/kube-no-trouble/pkg/judge"
	"github.com/doitintl/kube-no-trouble/pkg/printer"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var (
	version string = "dev"
	gitSha  string = "dev"
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
		log.Error().Err(err).Msg("Failed to initialize :" + collector.Name())
	} else {
		collectors = append(collectors, collector)
	}
	return collectors
}

func initCollectors(config *config.Config) []collector.Collector {
	collectors := []collector.Collector{}
	if config.Cluster {
		collector, err := collector.NewClusterCollector(&collector.ClusterOpts{Kubeconfig: config.Kubeconfig})
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

func main() {

	config, err := config.NewFromFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config flags")
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg(">>> Kube No Trouble `kubent` <<<")
	log.Info().Msgf("version %s (git sha %s)", version, gitSha)
	log.Info().Msg("Initializing collectors and retrieving data")

	initCollectors := initCollectors(config)
	collectors := getCollectors(initCollectors)

	judge, err := judge.NewRegoJudge(&judge.RegoOpts{})
	if err != nil {
		log.Fatal().Err(err).Str("name", "Rego").Msg("Failed to initialize decision engine")
	}

	results, err := judge.Eval(collectors)
	if err != nil {
		log.Fatal().Err(err).Str("name", "Rego").Msg("Failed to evaluate input")
	}

	var prnt printer.Printer
	switch config.Output {
	case "json":
		prnt, err = printer.NewJSONPrinter(&printer.JSONOpts{})
	default:
		prnt, err = printer.NewTextPrinter(&printer.TextOpts{})
	}
	if err != nil {
		log.Fatal().Err(err).Str("name", config.Output).Msg("Failed to initialize output printer")
	}

	err = prnt.Print(results)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to print results")
	}
}
