package printer

import (
	"github.com/stepanstipl/kube-no-trouble/pkg/judge"
)

type Printer interface {
	Print([]judge.Result) error
}
