package test

import (
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
)

type resourceFixtureTestCase struct {
	name          string
	fixturePaths  []string
	expectedKinds []string
}

func testReourcesUsingFixtures(t *testing.T, testCases []resourceFixtureTestCase) {
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

			for i := range manifests {
				if manifests[i]["kind"] != tc.expectedKinds[i] {
					t.Errorf("Expected to get %s, instead got: %s", tc.expectedKinds[i], manifests[i]["kind"])
				}
			}
		})
	}
}
