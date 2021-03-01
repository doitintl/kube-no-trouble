package rules

import (
	"embed"
	"path"
)

//go:embed rego
var local embed.FS

type Rule struct {
	Name string
	Rule string
}

func FetchRegoRules() ([]Rule, error) {
	fis, err := local.ReadDir("rego")
	if err != nil {
		return nil, err
	}

	rules := []Rule{}
	for _, info := range fis {
		data, err := local.ReadFile(path.Join("rego", info.Name()))
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{
			Name: info.Name(),
			Rule: string(data),
		})
	}

	return rules, nil
}
