package judge

type Result struct {
	Name        string
	Namespace   string
	Kind        string
	ApiVersion  string
	RuleSet     string
	ReplaceWith string
	Since       *Version
	Labels      map[string]interface{}
}

type Judge interface {
	Eval([]map[string]interface{}) ([]Result, error)
}
