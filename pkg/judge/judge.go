package judge

import "github.com/doitintl/kube-no-trouble/pkg/collector"

type Result struct {
	Name        string
	Namespace   string
	Kind        string
	ApiVersion  string
	RuleSet     string
	ReplaceWith string
	Since       *collector.Version
}

type Judge interface {
	Eval([]collector.MetaOject) ([]Result, error)
}
