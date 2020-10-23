package config

import (
	"path/filepath"

	flag "github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	Cluster    bool
	Debug      bool
	ExitError  bool
	Filenames  []string
	Helm2      bool
	Helm3      bool
	Kubeconfig string
	Output     string
}

func NewFromFlags() (*Config, error) {
	config := Config{}

	home := homedir.HomeDir()
	flag.BoolVarP(&config.Cluster, "cluster", "c", true, "enable Cluster collector")
	flag.BoolVarP(&config.Debug, "debug", "d", false, "enable debug logging")
	flag.BoolVarP(&config.ExitError, "exit-error", "e", false, "exit with non-zero code when issues are found")
	flag.BoolVar(&config.Helm2, "helm2", true, "enable Helm v2 collector")
	flag.BoolVar(&config.Helm3, "helm3", true, "enable Helm v3 collector")
	flag.StringSliceVarP(&config.Filenames, "filename", "f", []string{}, "manifests to check, use - for stdin")
	flag.StringVarP(&config.Kubeconfig, "kubeconfig", "k", filepath.Join(home, ".kube", "config"), "path to the kubeconfig file")
	flag.StringVarP(&config.Output, "output", "o", "text", "output format - [text|json]")

	flag.Parse()

	return &config, nil
}
