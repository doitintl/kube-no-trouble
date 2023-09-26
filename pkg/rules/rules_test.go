package rules

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestFetchRules(t *testing.T) {
	var expected []string
	root := "rego/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.Name() != "rego" {
			expected = append(expected, info.Name())
		}
		return nil
	})

	rules, err := FetchRegoRules([]schema.GroupVersionKind{})
	if err != nil {
		t.Errorf("Failed to load rules with: %s", err)
	}
	for i, rule := range rules {
		if rule.Name != expected[i] {
			t.Errorf("expected to get %s finding, instead got: %s", expected[i], rule.Name)
		}
	}
}

func TestFetchRulesWithAdditionalResources(t *testing.T) {
	var expected []string
	err := filepath.Walk(RULES_DIR, func(path string, info os.FileInfo, err error) error {
		if info.Name() != RULES_DIR {
			expected = append(expected, info.Name())
		}
		return nil
	})

	additionalKindsStr := []string{
		"ManagedCertificate.v1.networking.gke.io",
		"Fake.v1beta.example.com"}
	var additionalKinds []schema.GroupVersionKind
	for _, ar := range additionalKindsStr {
		gvr, _ := schema.ParseKindArg(ar)
		additionalKinds = append(additionalKinds, *gvr)
	}

	rules, err := FetchRegoRules([]schema.GroupVersionKind{})
	if err != nil {
		t.Errorf("Failed to load rules with: %s", err)
	}
	for i, rule := range rules {
		if rule.Name != expected[i] {
			t.Errorf("expected to get %s finding, instead got: %s", expected[i], rule.Name)
		}
	}
}

func TestRenderRuleRego(t *testing.T) {
	inputData := []byte("some input")
	fileName := "test.rego"

	outputData, err := renderRule(inputData, fileName, []schema.GroupVersionKind{})
	if err != nil {
		t.Errorf("Failed to render rule %s: %s", fileName, err)
	}

	if bytes.Compare(inputData, outputData) != 0 {
		t.Errorf("expected the input to be same as output")
	}
}

func TestRenderRuleTmpl(t *testing.T) {
	additionalResources := []schema.GroupVersionKind{
		schema.GroupVersionKind{
			Group:   "example.com",
			Version: "v2",
			Kind:    "Test",
		},
	}
	fileName := "test.tmpl"
	inputData := []byte("{{- range $key, $value := . }}" +
		"{{ $key }}: {{range $value}}{{.}}{{end}}" +
		"{{- end }}")
	expectedData := []byte("Test: example.com/v2")

	outputData, err := renderRule(inputData, fileName, additionalResources)
	if err != nil {
		t.Errorf("failed to render rule %s: %s", fileName, err)
	}

	if bytes.Compare(expectedData, outputData) != 0 {
		t.Errorf("result does not match expected output, expected: %s, got: %s", expectedData, outputData)
	}
}

func TestRenderRuleTmplFail(t *testing.T) {
	fileName := "test.tmpl"
	inputData := []byte("{{- rangeasd $key, $value := . }}" +
		"{{ $key }}: {{range $value}}{{.}}{{end}}" +
		"{{- end }}")

	_, err := renderRule(inputData, fileName, []schema.GroupVersionKind{})
	if err == nil {
		t.Errorf("expected this to fail")
	}
}

func TestRenderMultipleResources(t *testing.T) {
	additionalResources := []schema.GroupVersionKind{
		schema.GroupVersionKind{
			Group:   "example.com",
			Version: "v2",
			Kind:    "Test",
		},
		schema.GroupVersionKind{
			Group:   "example.com",
			Version: "v3",
			Kind:    "Test",
		},
		schema.GroupVersionKind{
			Group:   "example.com",
			Version: "v2",
			Kind:    "Result",
		},
	}
	fileName := "test.tmpl"
	inputData := []byte("{{- range $key, $value := . }}" +
		"{{ $key }}: {{range $value}}{{.}} {{end}}" +
		"{{- end }}")
	expectedData := []byte("Result: example.com/v2 Test: example.com/v2 example.com/v3")

	outputData, err := renderRule(inputData, fileName, additionalResources)
	if err != nil {
		t.Errorf("failed to render rule %s: %s", fileName, err)
	}

	outputData = []byte(strings.TrimSpace(string(outputData)))

	if bytes.Compare(expectedData, outputData) != 0 {
		t.Errorf("result does not match expected output, expected: %s, got: %s", expectedData, outputData)
	}
}

func TestRenderRuleUnknownFail(t *testing.T) {
	inputData := []byte("some input")
	fileName := "test.txt"

	_, err := renderRule(inputData, fileName, []schema.GroupVersionKind{})
	if err == nil {
		t.Errorf("expected this to fail")
	}
}
