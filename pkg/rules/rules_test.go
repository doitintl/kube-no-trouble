package rules

import (
	"os"
	"path/filepath"
	"testing"
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

	rules, err := FetchRegoRules()
	if err != nil {
		t.Errorf("Failed to load rules with: %s", err)
	}
	for i, rule := range rules {
		if rule.Name != expected[i] {
			t.Errorf("expected to get %s finding, instead got: %s", expected[i], rule.Name)
		}
	}
}
