package printer

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

func Test_newTextPrinter(t *testing.T) {
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
			got, err := newTextPrinter(tt.outputFileName)
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

func Test_textPrinter_Print(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tp := &textPrinter{
		commonPrinter: &commonPrinter{tmpFile},
	}

	version, _ := judge.NewVersion("1.2.3")
	results := []judge.Result{{"Name", "Namespace", "Kind", "1.2.3", "Test", "4.5.6", version}}

	if err := tp.Print(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fi, _ := tmpFile.Stat()
	if fi.Size() == 0 {
		t.Fatalf("expected non-zero size output file: %v", err)
	}
}

func Test_textPrinter_Close(t *testing.T) {
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
			c := &textPrinter{
				commonPrinter: &commonPrinter{tt.outputFile},
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
