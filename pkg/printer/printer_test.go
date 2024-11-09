package printer

import (
	"os"
	"reflect"
	"testing"
)

const (
	tempFilePrefix               = "kubent-tests-"
	tempFileCreateFailureMessage = "failed to create temp dir for testing: %v"
)

func TestNewPrinter(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name      string
		argChoice string
		options   PrinterOptions
		wantErr   bool
		want      Printer
	}{
		{"good", "json", PrinterOptions{outputFile: os.Stdout}, false, &jsonPrinter{&commonPrinter{}}},
		{"good-file", "json", PrinterOptions{outputFile: tmpFile}, false, &jsonPrinter{&commonPrinter{}}},
		{"invalid", "xxx", PrinterOptions{outputFile: tmpFile}, true, nil},
		{"empty", "", PrinterOptions{outputFile: os.Stdout}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrinter(tt.argChoice, &tt.options)
			if (err != nil) && !tt.wantErr {
				t.Errorf("unexpected error - got: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("unexpected result - got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_newCommonPrinter(t *testing.T) {
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
		{"good-stdout", PrinterOptions{outputFile: os.Stdout}, false},
		{"bad-empty", PrinterOptions{outputFile: nil, showLabels: false}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newCommonPrinter(&tt.options)
			if (err != nil) && !tt.wantErr {
				t.Fatalf("unexpected error - got: %v, wantErr: %v", err, tt.wantErr)
			}

			if err != nil && got != nil {
				t.Fatalf("expected nil return in case of an error, got %v", got)
			}

			if !tt.wantErr {
				defer got.Close()
				if (tt.options.outputFile.Name() != "-" && got.options.outputFile.Name() != tt.options.outputFile.Name()) ||
					tt.options.outputFile.Name() == "-" && got.options.outputFile.Name() != os.Stdout.Name() {
					t.Errorf("unexpected file name- got: %s, want: %s", got.options.outputFile.Name(), tt.options.outputFile.Name())
				}
			}
		})
	}
}

func Test_ensureOutputFileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}

	tests := []struct {
		name         string
		argsFileName string
		wantErr      bool
	}{
		{"stdout", "-", false},
		{"path", tmpFile.Name(), false},
		{"bad-dir", "/unlikely/to/exist/directory/my.log", true},
		{"bad-empty", "", true},
	}
	defer os.Remove(tmpFile.Name())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ensureOutputFileExists(tt.argsFileName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error - got: %v, wantErr: %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if got != nil {
					t.Fatalf("expected nil return in case of an error, got %v", got)
				}
			} else {
				if got == nil {
					t.Fatalf("unexpected nil return, got %v", got)
				}
				if got.Name() != tt.argsFileName && tt.argsFileName != "-" ||
					tt.argsFileName == "-" && got.Name() != os.Stdout.Name() {
					t.Fatalf("expected os.File with Name: %s, got: %s", tt.argsFileName, got.Name())
				}
			}
		})
	}
}

func Test_commonPrinter_Close(t *testing.T) {
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
		{"good-stdout", PrinterOptions{outputFile: os.Stdout}, false},
		{"bad-closed-file", PrinterOptions{outputFile: tmpFile}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commonPrinter{
				options: &tt.options,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
