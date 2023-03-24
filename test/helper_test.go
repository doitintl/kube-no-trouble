package test

import (
	"testing"

	"github.com/LeMyst/kube-no-trouble/pkg/collector"
	"github.com/LeMyst/kube-no-trouble/pkg/judge"
	"github.com/LeMyst/kube-no-trouble/pkg/rules"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type resourceFixtureTestCase struct {
	name          string
	fixturePaths  []string
	expectedKinds []string
}

func init() {
	Setup()
}

func testResourcesUsingFixtures(t *testing.T, testCases []resourceFixtureTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := collector.NewFileCollector(
				&collector.FileOpts{Filenames: tc.fixturePaths},
			)

			if err != nil {
				t.Errorf("Expected to succeed for %s, failed: %s", tc.fixturePaths, err)
			}

			manifests, err := c.Get()
			if err != nil {
				t.Errorf("Expected to succeed for %s, failed: %s", tc.fixturePaths, err)
			} else if len(manifests) != len(tc.expectedKinds) {
				t.Errorf("Expected to get %d, got %d", len(tc.expectedKinds), len(manifests))
			}

			loadedRules, err := rules.FetchRegoRules([]schema.GroupVersionKind{})
			if err != nil {
				t.Errorf("Failed to load rules")
			}

			judge, err := judge.NewRegoJudge(&judge.RegoOpts{}, loadedRules)
			if err != nil {
				t.Errorf("failed to create judge instance: %s", err)
			}

			results, err := judge.Eval(manifests)
			if err != nil {
				t.Errorf("failed to evaluate input: %s", err)
			}

			if len(results) != len(tc.expectedKinds) {
				t.Errorf("expected %d findings, instead got: %d", len(tc.expectedKinds), len(results))
			}

			for i := range manifests {
				if manifests[i]["kind"] != tc.expectedKinds[i] {
					t.Errorf("Expected to get %s, instead got: %s", tc.expectedKinds[i], manifests[i]["kind"])
				}
			}
		})
	}
}
