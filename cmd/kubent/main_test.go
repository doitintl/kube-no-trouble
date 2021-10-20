package main

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
	"github.com/doitintl/kube-no-trouble/pkg/config"
	"github.com/doitintl/kube-no-trouble/pkg/judge"

	"github.com/rs/zerolog"
)

func TestInitCollectors(t *testing.T) {
	testConfig := config.Config{
		Filenames:  []string{"../../fixtures/deployment-v1beta1.yaml"},
		Cluster:    false,
		Helm2:      false,
		Helm3:      false,
		Kubeconfig: "test",
		LogLevel:   config.ZeroLogLevel(zerolog.ErrorLevel),
		Output:     "text",
	}

	collectors := initCollectors(&testConfig)

	if collectors[0].Name() != "File" {
		t.Errorf("Did not parse fixture with path %s", testConfig.Filenames[0])
	}
}

func TestGetCollectors(t *testing.T) {
	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})

	if err != nil {
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	initCollectors := []collector.Collector{}
	initCollectors = append(initCollectors, fileCollector)

	collectors := getCollectors(initCollectors)

	if collectors != nil && len(collectors) != 1 {
		t.Errorf("Did not get file collector correctly with error: %s", err)
	}
}

func TestStoreCollector(t *testing.T) {
	collectors := []collector.Collector{}
	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})

	if err != nil {
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	collectors = storeCollector(fileCollector, err, collectors)

	if len(collectors) != 1 {
		t.Errorf("Failed to append collector")
	}
}

func TestStoreCollectorMultiple(t *testing.T) {
	collectors := []collector.Collector{}
	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})

	if err != nil {
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	collectors = storeCollector(fileCollector, err, collectors)

	collectors = storeCollector(fileCollector, err, collectors)

	if len(collectors) != 2 {
		t.Errorf("Failed to append collectors")
	}
}

func TestStoreCollectorError(t *testing.T) {
	collectors := []collector.Collector{}
	err := errors.New("Just testing...")

	collectors = storeCollector(nil, err, collectors)

	if len(collectors) != 0 {
		t.Errorf("Failed to ignore collector with error")
	}
}

func TestMainExitCodes(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string // file list
		expected int      // number of manifests
	}{
		{"success", []string{"-c=false", "--helm2=false", "--helm3=false"}, 0},
		{"errorBadFlag", []string{"-c=not-boolean"}, 2},
		{"successFound", []string{"-c=false", "--helm2=false", "--helm3=false", "-f=../../fixtures/deployment-v1beta1.yaml"}, 0},
		{"exitErrorFlagNone", []string{"-c=false", "--helm2=false", "--helm3=false", "-e"}, 0},
		{"exitErrorFlagFound", []string{"-c=false", "--helm2=false", "--helm3=false", "-e", "-f=../../fixtures/deployment-v1beta1.yaml"}, 200},
	}

	if os.Getenv("TEST_EXIT_CODE") == "1" {
		tc, err := strconv.Atoi(os.Getenv("TEST_CASE"))
		if err != nil {
			t.Errorf("failed to determin the test case num (TEST_CASE env var): %v", err)
		}

		os.Args = []string{os.Args[0]}
		os.Args = append(os.Args, testCases[tc].args...)
		main()
		return
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ee *exec.ExitError

			cmd := exec.Command(os.Args[0], "-test.run=TestMainExitCodes")
			cmd.Env = append(os.Environ(), "TEST_EXIT_CODE=1", "TEST_CASE="+strconv.Itoa(i))
			err := cmd.Run()

			if tc.expected == 0 && err != nil {
				t.Fatalf("expected to succeed with exit code %d, failed with %v", tc.expected, err)
			}

			if tc.expected != 0 && err == nil {
				t.Fatalf("expected to get exit code %d, succeeded with 0", tc.expected)
			}

			if tc.expected != 0 && err != nil {
				if errors.As(err, &ee) && ee.ExitCode() != tc.expected {
					t.Fatalf("expected to get exit code %d, failed with %v, exit code %d", tc.expected, err, ee.ExitCode())
				} else if !errors.As(err, &ee) {
					t.Fatalf("expected to get exit code %d, failed with %v", tc.expected, err)
				}
			}
		})
	}
}

func TestGetServerVersionNone(t *testing.T) {
	collectors := []collector.Collector{}

	version, err := getServerVersion(nil, collectors)
	if err != nil {
		t.Errorf("Failed to get version with error: %s", err)
	}

	if version != nil {
		t.Errorf("Expected no version to be detected, instead got: %s", version.String())
	}
}

func TestGetServerVersionNotSupported(t *testing.T) {
	collectors := []collector.Collector{}

	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})
	if err != nil {
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	collectors = storeCollector(fileCollector, err, collectors)
	version, err := getServerVersion(nil, collectors)
	if err != nil {
		t.Errorf("Failed to get version with error: %s", err)
	}

	if version != nil {
		t.Errorf("Expected no version to be detected, instead got: %s", version.String())
	}
}

func TestGetServerVersion(t *testing.T) {
	collectors := []collector.Collector{}

	fc := collector.NewFakeCollector()
	collectors = storeCollector(fc, nil, collectors)

	version, err := getServerVersion(nil, collectors)
	if err != nil {
		t.Errorf("Failed to get version with error: %s", err)
	}

	fakeVersion, err := judge.NewVersion(collector.FAKE_VERSION)
	if version.Compare(fakeVersion.Version) != 0 {
		t.Errorf("Expected %s version to be detected, instead got: %s", fakeVersion.String(), version.String())
	}
}
