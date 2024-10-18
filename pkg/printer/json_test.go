package printer

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

func Test_newJSONPrinter(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name    string
		options PrinterOptions
		wantErr bool
	}{
		{"good-stdout", PrinterOptions{outputFile: os.Stdout}, false},
		{"good-file", PrinterOptions{outputFile: tmpFile}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newJSONPrinter(&tt.options)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v, expected error: %v", err, tt.wantErr)
			}

			if err != nil || got == nil {
				t.Errorf("expected nil in case of an error, got %v", got)
			}
		})
	}
}
func Test_jsonPrinter_Print(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	options := &PrinterOptions{outputFile: tmpFile}
	c := &jsonPrinter{
		commonPrinter: &commonPrinter{options},
	}

	results := getTestResult(map[string]interface{}{"key2": "value2"})

	if err := c.Print(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpFile.Seek(0, 0)

	var readResults []judge.Result
	readBytes, err := os.ReadFile(tmpFile.Name())
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
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name    string
		options PrinterOptions
		wantErr bool
	}{
		{"good-file", PrinterOptions{outputFile: tmpFile}, false},
		{"bad-closed-file", PrinterOptions{outputFile: tmpFile}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &jsonPrinter{
				commonPrinter: &commonPrinter{&tt.options},
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
