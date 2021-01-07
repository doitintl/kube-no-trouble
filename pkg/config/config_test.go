package config

import (
	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
	"os"
	"testing"
)

func TestNewFromFlags(t *testing.T) {
	// reset for testing
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	config, err := NewFromFlags()

	if err != nil {
		t.Errorf("Flags parsing failed %s", err)
	}

	if !config.Cluster && config.Output != "text" {
		t.Errorf("Config not parsed correctly")
	}
}

func TestNewFromFlagsKubeconfigEnv(t *testing.T) {
	// reset for testing
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	testVal := "my-file.conf"

	err := os.Setenv("KUBECONFIG", testVal)
	if err != nil {
		t.Errorf("failed to set env variable: %s", err)
	}

	config, err := NewFromFlags()

	if err != nil {
		t.Errorf("Formatting flags failed %s", err)
	}

	if config.Kubeconfig != testVal {
		t.Errorf("kubeconfig option not loaded correctly from ebv variable, expected: %s, got: %s", testVal, config.Kubeconfig)
	}
}

func TestNewFromFlagsKubeconfigHome(t *testing.T) {
	// reset for testing
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	expected := homedir.HomeDir() + "/.kube/config"

	err := os.Unsetenv("KUBECONFIG")
	if err != nil {
		t.Errorf("failed to unset KUBECONFIG env variable: %s", err)
	}

	config, err := NewFromFlags()

	if err != nil {
		t.Errorf("Formatting flags failed %s", err)
	}

	if config.Kubeconfig != expected {
		t.Errorf("kubeconfig option not set to expected default, expected: %s, got: %s", expected, config.Kubeconfig)
	}
}

func TestEnvOrStringVariable(t *testing.T) {
	err := os.Setenv("FOO", "1")
	if err != nil {
		t.Errorf("failed to set env variable: %e", err)
	}

	i := envOrString("FOO", "default")
	if i != "1" {
		t.Errorf("Expected to get env variable, got %s insteadt", i)
	}
}

func TestEnvOrStringDefault(t *testing.T) {
	err := os.Unsetenv("FOO")
	if err != nil {
		t.Errorf("failed to unset env variable: %e", err)
	}

	i := envOrString("FOO", "default")
	if i != "default" {
		t.Errorf("Expected to get default string, got %s instead", i)
	}
}
