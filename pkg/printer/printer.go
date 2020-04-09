package printer

import (
	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type Printer interface {
	Print([]judge.Result) error
}
