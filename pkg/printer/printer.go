package printer

import (
	"fmt"
	"os"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

type Printer interface {
	Print([]judge.Result) error
	Close() error
}

type PrinterOptions struct {
	showLabels bool
	outputFile *os.File
}

type commonPrinter struct {
	options *PrinterOptions
}

func NewPrinterOptions(fileName string, showLabels bool) (*PrinterOptions, error) {
	outputFile, err := ensureOutputFileExists(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure output device: %v", err)
	}
	return &PrinterOptions{
		showLabels,
		outputFile,
	}, nil
}

// newCommonPrinter creates new printer that prints to given output file
func newCommonPrinter(options *PrinterOptions) (*commonPrinter, error) {
	return &commonPrinter{
		options,
	}, nil
}

// NewPrinter creates new printer of given type that prints to given output file
func NewPrinter(choice string, options *PrinterOptions) (Printer, error) {
	switch choice {
	case "json":
		return newJSONPrinter(options)
	case "text":
		return newTextPrinter(options)
	case "csv":
		return newCSVPrinter(options)
	default:
		return nil, fmt.Errorf("unknown printer type %s", choice)
	}
}

// Close will free resources used by the printer
func (c *commonPrinter) Close() error {
	if c.options.outputFile.Name() != os.Stdout.Name() {
		if err := c.options.outputFile.Close(); err != nil {
			return fmt.Errorf("failed to close output file: %w", err)
		}
	}

	return nil
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
