package judge

type Result struct {
	Name       string
	Namespace  string
	Kind       string
	ApiVersion string
	RuleSet    string
}

type Judge interface {
	Eval([]map[string]interface{}) ([]Result, error)
}
