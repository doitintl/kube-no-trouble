package rules

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const RULES_DIR = "rego"

//go:embed rego
var local embed.FS

type Rule struct {
	Name string
	Rule string
}

func FetchRegoRules(additionalResources []schema.GroupVersionKind) ([]Rule, error) {
	fis, err := local.ReadDir(RULES_DIR)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	for _, info := range fis {
		data, err := local.ReadFile(path.Join(RULES_DIR, info.Name()))
		if err != nil {
			return nil, err
		}

		rule, err := renderRule(data, info.Name(), additionalResources)
		if err != nil {
			return nil, err
		}

		rules = append(rules, Rule{
			Name: info.Name(),
			Rule: string(rule),
		})
	}

	return rules, nil
}

func renderRule(inputData []byte, fileName string, additionalKinds []schema.GroupVersionKind) ([]byte, error) {
	var data []byte

	switch {
	case strings.HasSuffix(fileName, ".rego"):
		data = inputData

	// currently this is relevant only to additional resources
	case strings.HasSuffix(fileName, ".tmpl"):
		t, err := template.New(fileName).Parse(string(inputData))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", fileName, err)
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, additionalKinds); err != nil {
			return nil, fmt.Errorf("failed to render template %s: %w", fileName, err)
		}

		data = tpl.Bytes()

	default:
		return nil, fmt.Errorf("unrecognized filetype: %s", fileName)
	}

	return data, nil
}
