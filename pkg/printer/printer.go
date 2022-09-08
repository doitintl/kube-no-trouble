package printer

import (
	"fmt"
	"os"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

var printers = map[string]func(string) (Printer, error){
	"json": newJSONPrinter,
	"text": newTextPrinter,
}

type Printer interface {
	Print([]judge.Result) error
	Close() error
}

type commonPrinter struct {
	outputFile *os.File
}

// newCommonPrinter creates new printer that prints to given output file
func newCommonPrinter(outputFileName string) (*commonPrinter, error) {
	outputFile, err := ensureOutputFileExists(outputFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure output device: %v", err)
	}

	return &commonPrinter{
		outputFile: outputFile,
	}, nil
}

// NewPrinter creates new printer of given type that prints to given output file
func NewPrinter(choice string, outputFileName string) (Printer, error) {
	printer, err := ParsePrinter(choice)
	if err != nil {
		return nil, err
	}
	return printer(outputFileName)
}

// Close will free resources used by the printer
func (c *commonPrinter) Close() error {
	if c.outputFile.Name() != os.Stdout.Name() {
		if err := c.outputFile.Close(); err != nil {
			return fmt.Errorf("failed to close output file: %w", err)
		}
	}

	return nil
}

func ParsePrinter(choice string) (func(string) (Printer, error), error) {
	printer, exists := printers[choice]
	if !exists {
		return nil, fmt.Errorf("unknown printer type %s", choice)
	}
	return printer, nil
}

// ensureOutputFileExists will open file for writing, or create one if non-existent
func ensureOutputFileExists(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdout, nil
	}

	of, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create/open output file %s: %w", fileName, err)
	}

	return of, nil
}
