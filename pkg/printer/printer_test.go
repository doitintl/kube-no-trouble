package printer

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const (
	tempFilePrefix               = "kubent-tests-"
	tempFileCreateFailureMessage = "failed to create temp dir for testing: %v"
)

func TestParsePrinter(t *testing.T) {
	for k, v := range printers {
		p, err := ParsePrinter(k)
		if err != nil {
			t.Fatalf("failed to parse printer %s: %v", k, err)
		}

		if reflect.ValueOf(p).Pointer() != reflect.ValueOf(v).Pointer() {
			t.Fatalf("expected to get function %p, got %p instead", p, p)
		}
	}
}

func TestParsePrinterInvalid(t *testing.T) {
	_, err := ParsePrinter("BAD")
	if err == nil {
		t.Fatalf("expected ParsePrinter to fail with unimplemented type")
	}
}

func TestNewPrinter(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name              string
		argChoice         string
		argOutputFileName string
		wantErr           bool
		want              Printer
	}{
		{"good", "json", "-", false, &jsonPrinter{&commonPrinter{}}},
		{"good-file", "json", tmpFile.Name(), false, &jsonPrinter{&commonPrinter{}}},
		{"invalid", "xxx", "", true, nil},
		{"invalid-out", "json", "/not/likely/to/exist", true, nil},
		{"empty", "", "-", true, nil},
		{"empty-out", "json", "", true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrinter(tt.argChoice, tt.argOutputFileName)
			if (err != nil) != tt.wantErr {
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
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		t.Fatalf(tempFileCreateFailureMessage, err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name              string
		argOutputFileName string
		wantErr           bool
	}{
		{"good-file", tmpFile.Name(), false},
		{"good-stdout", "-", false},
		{"bad-empty", "", true},
		{"bad-path", "/this/is/unlikely/to/exist", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newCommonPrinter(tt.argOutputFileName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error - got: %v, wantErr: %v", err, tt.wantErr)
			}

			if err != nil && got != nil {
				t.Fatalf("expected nil return in case of an error, got %v", got)
			}

			if !tt.wantErr {
				defer got.Close()
				if (tt.argOutputFileName != "-" && got.outputFile.Name() != tt.argOutputFileName) ||
					tt.argOutputFileName == "-" && got.outputFile.Name() != os.Stdout.Name() {
					t.Errorf("unexpected file name- got: %s, want: %s", got.outputFile.Name(), tt.argOutputFileName)
				}
			}
		})
	}
}

func Test_ensureOutputFileExists(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
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
		{"good-stdout", os.Stdout, false},
		{"bad-closed-file", tmpFile, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commonPrinter{
				outputFile: tt.outputFile,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error - got: %v, expected error: %v", err, tt.wantErr)
			}
		})
	}
}
