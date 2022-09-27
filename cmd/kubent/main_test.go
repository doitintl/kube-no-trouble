package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
	"github.com/doitintl/kube-no-trouble/pkg/config"
	"github.com/doitintl/kube-no-trouble/pkg/judge"

	"github.com/rs/zerolog"
)

const FIXTURES_DIR = "../../fixtures"

func TestInitCollectors(t *testing.T) {
	testConfig := config.Config{
		Filenames:  []string{filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")},
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
	fileCollector, err := collector.NewFileCollector(
		&collector.FileOpts{Filenames: []string{filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}})

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
	fileCollector, err := collector.NewFileCollector(
		&collector.FileOpts{Filenames: []string{filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}})

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
	fileCollector, err := collector.NewFileCollector(
		&collector.FileOpts{Filenames: []string{filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}})

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
	tmpDir, err := ioutil.TempDir(os.TempDir(), "kubent-tests-")
	if err != nil {
		t.Fatalf("failed to create temp dir for testing: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	expectedJsonOutput, _ := os.ReadFile(filepath.Join(FIXTURES_DIR, "expected-json-output.json"))
	helm2FlagDisabled := "--helm2=false"
	helm3FlagDisabled := "--helm3=false"
	clusterFlagDisabled := "--cluster=false"
	testCases := []struct {
		name        string
		args        []string // file list
		expected    int      // expected exit code
		stdout      string   // expected stdout
		outFileName string
	}{
		{"success", []string{clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled}, 0, "", ""},
		{"errorBadFlag", []string{"-c=not-boolean"}, 2, "", ""},
		{"successFound", []string{"-o=json", clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, string(expectedJsonOutput), ""},
		{"exitErrorFlagNone", []string{clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled, "-e"}, 0, "", ""},
		{"exitErrorFlagFound", []string{clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled, "-e", "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 200, "", ""},
		{"version short flag set", []string{"-v"}, 0, "", ""},
		{"version long flag set", []string{"--version"}, 0, "", ""},
		{"empty text output", []string{clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled}, 0, "", ""},
		{"empty json output", []string{"-o=json", clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled}, 0, "[]\n", ""},
		{"json-file", []string{"-o=json", clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, "", filepath.Join(tmpDir, "json-file.out")},
		{"text-file", []string{"-o=json", clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, "", filepath.Join(tmpDir, "text-file.out")},
		{"json-stdout", []string{"-o=json", clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, string(expectedJsonOutput), "-"},
		{"error-bad-file", []string{clusterFlagDisabled, helm2FlagDisabled, helm3FlagDisabled}, 1, "", "/this/dir/is/unlikely/to/exist"},
	}

	if os.Getenv("TEST_EXIT_CODE") == "1" {
		args := []string{}
		decodeBase64(&args, os.Getenv("TEST_ARGS"))
		os.Args = append(os.Args[:1], args...)
		main()
		return
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ee *exec.ExitError

			if tc.outFileName != "" {
				tc.args = append(tc.args, "-O="+tc.outFileName)
			}
			base64Args, _ := encodeBase64(tc.args)

			cmd := exec.Command(os.Args[0], "-test.run=TestMainExitCodes")
			cmd.Env = append(os.Environ(),
				"TEST_EXIT_CODE=1",
				"TEST_ARGS="+base64Args)
			out, err := cmd.Output()

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
			if tc.expected == 0 && err == nil && tc.stdout != string(out) {
				t.Fatalf("expected to get stdout as %s, instead got %s", tc.stdout, out)
			}

			if tc.expected == 0 && err == nil && tc.outFileName != "" && tc.outFileName != "-" {
				if fs, err := os.Stat(tc.outFileName); err != nil || fs.Size() == 0 {
					t.Fatalf("expected non-empty outputdile: %s, got error: %v", tc.outFileName, err)
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

	fileCollector, err := collector.NewFileCollector(
		&collector.FileOpts{Filenames: []string{filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}})
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

	fakeVersion, _ := collector.NewVersion(collector.FAKE_VERSION)
	if version.Compare(fakeVersion.Version) != 0 {
		t.Errorf("Expected %s version to be detected, instead got: %s", fakeVersion.String(), version.String())
	}
}

func encodeBase64(args []string) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(args)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

func decodeBase64(dst *[]string, encoded string) error {
	r := strings.NewReader(encoded)
	base64Dec := base64.NewDecoder(base64.StdEncoding, r)
	jsonDec := json.NewDecoder(base64Dec)

	return jsonDec.Decode(dst)
}

func Test_outputResults(t *testing.T) {
	testVersion, _ := collector.NewVersion("4.5.6")
	testResults := []judge.Result{{"name", "ns", "kind",
		"1.2.3", "rs", "rep", testVersion}}

	type args struct {
		results    []judge.Result
		outputType string
		outputFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"good", args{testResults, "text", "-"}, false},
		{"bad-new-printer-type", args{testResults, "unknown", "-"}, true},
		{"bad-new-printer-file", args{testResults, "text", "/unlikely/to/exist/dir"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResults(tt.args.results, tt.args.outputType, tt.args.outputFile); (err != nil) != tt.wantErr {
				t.Errorf("unexpected error - got: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
