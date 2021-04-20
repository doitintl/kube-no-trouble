package printer

import (
	"fmt"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

var printers = map[string]func() Printer{
	"json": newJSONPrinter,
	"text": newTextPrinter,
}

type Printer interface {
	Print([]judge.Result) error
}

func NewPrinter(choice string) (Printer, error) {
	printer, err := ParsePrinter(choice)
	if err != nil {
		return nil, err
	}
	return printer(), nil
}

func ParsePrinter(choice string) (func() Printer, error) {
	printer, exists := printers[choice]
	if !exists {
		return nil, fmt.Errorf("unknown printer type %s", choice)
	}
	return printer, nil
}
