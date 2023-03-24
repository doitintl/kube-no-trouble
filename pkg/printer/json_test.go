package printer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

func Test_newJSONPrinter(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name           string
		outputFileName string
		wantErr        bool
	}{
		{"good-stdout", "-", false},
		{"good-file", tmpFile.Name(), false},
		{"bad-empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newJSONPrinter(tt.outputFileName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v, expected error: %v", err, tt.wantErr)
			}

			if err != nil && got != nil {
				t.Errorf("expected nil in case of an error, got %v", got)
			}
		})
	}
}

func Test_jsonPrinter_Print(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	c := &jsonPrinter{
		commonPrinter: &commonPrinter{tmpFile},
	}

	version, _ := judge.NewVersion("1.2.3")
	results := []judge.Result{{"Name", "Namespace", "Kind", "1.2.3", "Test", "4.5.6", version}}

	if err := c.Print(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpFile.Seek(0, 0)

	var readResults []judge.Result
	readBytes, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error reading back the file: %v", err)
	}
	if err := json.Unmarshal(readBytes, &readResults); err != nil {
		t.Fatalf("unexpected error unmarshalling the previously written file: %v", err)
	}
	if !reflect.DeepEqual(readResults, results) {
		t.Fatalf("written and read result do not seem to be equal")
	}
}

func Test_jsonPrinter_Close(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name       string
		outputFile *os.File
		wantErr    bool
	}{
		{"good-file", tmpFile, false},
		{"bad-closed-file", tmpFile, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &jsonPrinter{
				commonPrinter: &commonPrinter{tt.outputFile},
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
