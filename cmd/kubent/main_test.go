package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
	"github.com/doitintl/kube-no-trouble/pkg/config"
	ctxKey "github.com/doitintl/kube-no-trouble/pkg/context"
	"github.com/doitintl/kube-no-trouble/pkg/judge"

	"github.com/rs/zerolog"
)

const FIXTURES_DIR = "../../fixtures"

func TestInitCollectors(t *testing.T) {
	testConfig := config.Config{
		Filenames:  []string{filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")},
		Cluster:    false,
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

	collectors, _ := getCollectors(initCollectors)

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
	tmpDir, err := os.MkdirTemp(os.TempDir(), "kubent-tests-")
	if err != nil {
		t.Fatalf("failed to create temp dir for testing: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	expectedJsonOutput, _ := os.ReadFile(filepath.Join(FIXTURES_DIR, "expected-json-output.json"))
	expectedJsonOutputLabels, _ := os.ReadFile(filepath.Join(FIXTURES_DIR, "expected-json-output-labels.json"))
	helm3FlagDisabled := "--helm3=false"
	clusterFlagDisabled := "--cluster=false"
	testCases := []struct {
		name        string
		args        []string // file list
		expected    int      // expected exit code
		stdout      string   // expected stdout
		outFileName string
		emptyStderr bool
	}{
		{"success", []string{clusterFlagDisabled, helm3FlagDisabled}, 0, "", "", false},
		{"errorBadFlag", []string{"-c=not-boolean"}, 2, "", "", false},
		{"successFound", []string{"-o=json", clusterFlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, string(expectedJsonOutput), "", false},
		{"successFoundWithLabels", []string{"--labels=true", "-o=json", clusterFlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1-labels.yaml")}, 0, string(expectedJsonOutputLabels), "", false},
		{"exitErrorFlagNone", []string{clusterFlagDisabled, helm3FlagDisabled, "-e"}, 0, "", "", false},
		{"exitErrorFlagFound", []string{clusterFlagDisabled, helm3FlagDisabled, "-e", "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 200, "", "", false},
		{"version short flag set", []string{"-v"}, 0, "", "", false},
		{"version long flag set", []string{"--version"}, 0, "", "", false},
		{"empty text output", []string{clusterFlagDisabled, helm3FlagDisabled}, 0, "", "", false},
		{"empty json output", []string{"-o=json", clusterFlagDisabled, helm3FlagDisabled}, 0, "[]\n", "", false},
		{"fail to get collectors", []string{"-o=json", "-f=fail"}, 100, "[]\n", "", false},
		{"json-file", []string{"-o=json", clusterFlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, "", filepath.Join(tmpDir, "json-file.out"), false},
		{"text-file", []string{"-o=text", clusterFlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, "", filepath.Join(tmpDir, "text-file.out"), false},
		{"json-stdout", []string{"-o=json", clusterFlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1.yaml")}, 0, string(expectedJsonOutput), "-", false},
		{"json-stdout-with-labels", []string{"--labels=true", "-o=json", clusterFlagDisabled, helm3FlagDisabled, "-f=" + filepath.Join(FIXTURES_DIR, "deployment-v1beta1-labels.yaml")}, 0, string(expectedJsonOutputLabels), "-", false},
		{"error-bad-file", []string{clusterFlagDisabled, helm3FlagDisabled}, 1, "", "/this/dir/is/unlikely/to/exist", false},
		{"no-3rdparty-output", []string{clusterFlagDisabled, helm3FlagDisabled, "-l=disabled"}, 0, "", "", true},
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

			var stdout, stderr bytes.Buffer

			cmd := exec.Command(os.Args[0], "-test.run=TestMainExitCodes")
			cmd.Env = append(os.Environ(),
				"TEST_EXIT_CODE=1",
				"TEST_ARGS="+base64Args)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()

			outStr := stdout.String()
			errStr := stderr.String()

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
			if tc.expected == 0 && err == nil && tc.stdout != outStr {
				t.Fatalf("expected to get stdout as %s, instead got %s", tc.stdout, outStr)
			}

			if tc.expected == 0 && err == nil && tc.outFileName != "" && tc.outFileName != "-" {
				if fs, err := os.Stat(tc.outFileName); err != nil || fs.Size() == 0 {
					t.Fatalf("expected non-empty outputfile: %s, got error: %v", tc.outFileName, err)
				}
			}

			if tc.emptyStderr && errStr != "" {
				t.Fatalf("expected empty stderr, got: %s", errStr)
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

	fakeVersion, _ := judge.NewVersion(collector.FAKE_VERSION)
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
	testVersion, _ := judge.NewVersion("4.5.6")
	testResults := []judge.Result{{Name: "name", Namespace: "ns", Kind: "kind",
		ApiVersion: "1.2.3", RuleSet: "rs", ReplaceWith: "rep", Since: testVersion}}

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

	labelsFlag := false
	ctx := context.WithValue(context.Background(), ctxKey.LABELS_CTX_KEY, &labelsFlag)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResults(tt.args.results, tt.args.outputType, tt.args.outputFile, ctx); (err != nil) != tt.wantErr {
				t.Errorf("unexpected error - got: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configureGlobalLogging(t *testing.T) {
	// just make sure the method runs, this is mostly covered
	//by the Test_MainExitCodes
	configureGlobalLogging()
}
