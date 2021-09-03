package judge

import (
	goversion "github.com/hashicorp/go-version"
)

type Result struct {
	Name        string
	Namespace   string
	Kind        string
	ApiVersion  string
	RuleSet     string
	ReplaceWith string
	Since       *goversion.Version `json:"-"`
	SinceStr    string             `json:"Since"`
}

type Judge interface {
	Eval([]map[string]interface{}) ([]Result, error)
}
