package config

import (
	"context"
	"os"
	"testing"

	goversion "github.com/hashicorp/go-version"

	"github.com/spf13/pflag"
)

func TestValidLogLevelFromFlags(t *testing.T) {
	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()
	ctx := context.Background()

	var validLevels = []string{"trace", "debug", "info", "warn", "error", "fatal", "panic", "", "disabled"}
	for i, level := range validLevels {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--log-level=" + level

		config, _, err := NewFromFlags(ctx)

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
	ctx := context.Background()

	config, _, err := NewFromFlags(ctx)

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
		{"abcdef"},
		{""},
		{"test.v1.com"},
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
	ctx := context.Background()

	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	for _, v := range validVersions {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--target-version=" + v
		config, _, err := NewFromFlags(ctx)

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
	ctx := context.Background()
	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	for _, v := range invalidVersions {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

		os.Args[1] = "--target-version=" + v
		config, _, _ := NewFromFlags(ctx)

		if config.TargetVersion != nil {
			t.Errorf("expected --target-version flag parsing to fail for: %s", v)
		}
	}
}

func TestContext(t *testing.T) {
	validContexts := []string{
		"my-context",
	}
	ctx := context.Background()
	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	for _, context := range validContexts {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--context=" + context
		config, _, err := NewFromFlags(ctx)

		if err != nil {
			t.Errorf("Flags parsing failed %s", err)
		}

		if config.Context != context {
			t.Fatalf("Context not parsed correctly: expected: %s, got: %s", context, config.Context)
		}
	}
}

func Test_validateOutputFile(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		wantErr bool
	}{
		{"empty", "", true},
		{"does-not-exist", "/this/directory/is/unlikely/to/exist", true},
		{"relative", "my.log", false},
		{"absolute", "/my.log", false},
		{"stdout", "-", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOutputFile(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("expected error = %v, got %v instead", err, tt.wantErr)
			}
		})
	}
}
