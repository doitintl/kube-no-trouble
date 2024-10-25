package printer

import (
	"os"
	"testing"
)

func Test_newTextPrinter(t *testing.T) {
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
			got, err := newTextPrinter(&tt.options)
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
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())
	options := &PrinterOptions{outputFile: tmpFile}
	tp := &textPrinter{
		commonPrinter: &commonPrinter{options},
	}

	results := getTestResult(map[string]interface{}{"key2": "value2"})

	if err := tp.Print(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fi, _ := tmpFile.Stat()
	if fi.Size() == 0 {
		t.Fatalf("expected non-zero size output file: %v", err)
	}
}

func Test_textPrinter_Close(t *testing.T) {
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
			c := &textPrinter{
				commonPrinter: &commonPrinter{&tt.options},
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
