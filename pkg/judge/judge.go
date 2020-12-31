package judge

type Result struct {
	Name        string
	Namespace   string
	Kind        string
	ApiVersion  string
	RuleSet     string
	ReplaceWith string
	Since       string
}

type Judge interface {
	Eval([]map[string]interface{}) ([]Result, error)
}
