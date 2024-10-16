package judge

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/rules"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const FIXTURES_DIR = "../../fixtures"

func TestNewRegoJudge(t *testing.T) {
	_, err := NewRegoJudge(&RegoOpts{}, []rules.Rule{})
	if err != nil {
		t.Errorf("failed to create judge instance: %s", err)
	}
}

func TestEvalEmpty(t *testing.T) {
	inputs := []map[string]interface{}{}

	judge, err := NewRegoJudge(&RegoOpts{}, []rules.Rule{})
	if err != nil {
		t.Errorf("failed to create judge instance: %s", err)
	}

	results, err := judge.Eval(inputs)
	if err != nil {
		t.Errorf("failed to evaluate input: %s", err)
	}

	if results == nil || len(results) != 0 {
		t.Errorf("expected empty array, instead got: %v", results)
	}
}

func TestEvalRules(t *testing.T) {
	testCases := []struct {
		name       string
		inputFiles []string // file list
		expected   []string // findings - kinds
	}{
		{"deprecated-1-16.rego", []string{"deployment-v1beta1.yaml"}, []string{"Deployment"}},
		{"deprecated-1-22.rego", []string{"ingress-v1beta1.yaml"}, []string{"Ingress"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var manifests []map[string]interface{}

			for _, f := range tc.inputFiles {
				var input []byte
				var err error

				input, err = ioutil.ReadFile(filepath.Join(FIXTURES_DIR, f))
				if err != nil {
					t.Errorf("failed to read file %s: %v", f, err)
				}

				var manifest map[string]interface{}
				err = yaml.Unmarshal([]byte(input), &manifest)
				if err != nil {
					t.Errorf("failed to parse file %s: %v", f, err)
				}

				manifests = append(manifests, manifest)
			}

			loadedRules, err := rules.FetchRegoRules([]schema.GroupVersionKind{})
			if err != nil {
				t.Errorf("Failed to load rules")
			}

			judge, err := NewRegoJudge(&RegoOpts{}, loadedRules)
			if err != nil {
				t.Errorf("failed to create judge instance: %s", err)
			}

			results, err := judge.Eval(manifests)
			if err != nil {
				t.Errorf("failed to evaluate input: %s", err)
			}

			if len(results) != len(tc.expected) {
				t.Errorf("expected %d findings, instead got: %d", len(tc.expected), len(results))
			}

			for i := range results {
				if results[i].Kind != tc.expected[i] {
					t.Errorf("expected to get %s finding, instead got: %s", tc.expected[i], results[i].Kind)
				}
			}
		})
	}

}
