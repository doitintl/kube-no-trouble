package printer

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	ctxKey "github.com/LeMyst/kube-no-trouble/pkg/context"
)

func TestNewCSVPrinter(t *testing.T) {
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
			got, err := newCSVPrinter(tt.outputFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && got != nil {
				t.Errorf("expected nil in case of an error, got %v", got)
			}
		})
	}
}

func TestCSVPrinterPrint(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tp := &csvPrinter{
		commonPrinter: &commonPrinter{tmpFile},
	}

	labelsFlag := false
	ctx := context.WithValue(context.Background(), ctxKey.LABELS_CTX_KEY, &labelsFlag)

	results := getTestResult(map[string]interface{}{"key2": "value2"})

	if err := tp.Print(results, ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fi, _ := tmpFile.Stat()
	if fi.Size() == 0 {
		t.Fatalf("expected non-zero size output file: %v", err)
	}
}

func TestCSVPrinterClose(t *testing.T) {
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
			c := &csvPrinter{
				commonPrinter: &commonPrinter{tt.outputFile},
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
