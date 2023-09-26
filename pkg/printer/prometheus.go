package printer

import (
	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type PrometheusPrinter struct {
}

func newPrometheusPrinter(outputFileName string) (Printer, error) {
	return &PrometheusPrinter{}, nil
}

func (p *PrometheusPrinter) Print([]judge.Result) error {
	return nil
}

func (p *PrometheusPrinter) Close() error {
	return nil
}
