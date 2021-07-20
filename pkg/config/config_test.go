package config

import (
	goversion "github.com/hashicorp/go-version"
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func TestValidLogLevelFromFlags(t *testing.T) {
	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	var validLevels = []string{"trace", "debug", "info", "warn", "error", "fatal", "panic", "", "disabled"}
	for i, level := range validLevels {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--log-level=" + level
		config, err := NewFromFlags()

		if err != nil {
			t.Errorf("Flags parsing failed %s", err)
		}

		expected := ZeroLogLevel(i - 1)
		actual := config.LogLevel

		if actual != expected {
			t.Errorf("Config not parsed correctly: %s \nactual %d, expected %d", level, actual, expected)
		}
	}
}

func TestInvalidLogLevelFromFlags(t *testing.T) {
	var testLevel ZeroLogLevel

	if err := testLevel.Set("bad"); err == nil {
		t.Errorf("Should not parse invalid flag")
	}
}

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

func TestValidateAdditionalResources(t *testing.T) {
	resources := []string{
		"Test.v1.example.com",
		"ManagedCertificates.v1.networking.gke.io",
		"ManagedCertificates.networking.gke.io",
	}

	err := validateAdditionalResources(resources)

	if err != nil {
		t.Errorf("expected resources %s to pass validation: %s", resources, err)
	}
}

func TestValidateAdditionalResourcesFail(t *testing.T) {
	testCases := [][]string{
		[]string{"abcdef"},
		[]string{""},
		[]string{"test.v1.com"},
	}

	for _, tc := range testCases {
		err := validateAdditionalResources(tc)

		if err == nil {
			t.Errorf("expected resources %s to fail validation: %s", tc, err)
		}
	}
}

func TestTargetVersion(t *testing.T) {
	validVersions := []string{
		"1.16", "1.16.3", "1.2.3",
	}

	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	for _, v := range validVersions {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--target-version=" + v
		config, err := NewFromFlags()

		if err != nil {
			t.Errorf("Flags parsing failed %s", err)
		}

		expected, _ := goversion.NewVersion(v)
		if config.TargetVersion.Version == nil {
			t.Fatalf("Target version not parsed correctly: expected: %s, got: %s", expected.String(), config.TargetVersion)
		}

		if !config.TargetVersion.Equal(expected) {
			t.Fatalf("Target version not parsed correctly: expected: %s, got: %s", expected.String(), config.TargetVersion.String())
		}
	}
}

func TestTargetVersionInvalid(t *testing.T) {
	invalidVersions := []string{
		"1.blah", "nope",
	}

	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	for _, v := range invalidVersions {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

		os.Args[1] = "--target-version=" + v
		config, _ := NewFromFlags()

		if config.TargetVersion.Version != nil {
			t.Errorf("expected --target-version flag parsing to fail for: %s", v)
		}
	}
}

func TestContext(t *testing.T) {
	validContexts := []string{
		"my-context",
	}

	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	for _, context := range validContexts {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--context=" + context
		config, err := NewFromFlags()

		if err != nil {
			t.Errorf("Flags parsing failed %s", err)
		}

		if config.Context != context {
			t.Fatalf("Context not parsed correctly: expected: %s, got: %s", context, config.Context)
		}
	}
}
