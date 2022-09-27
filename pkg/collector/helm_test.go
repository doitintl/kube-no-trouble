package collector

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSplitManifests(t *testing.T) {
	testCases := []struct {
		name     string
		input    string // manifest
		expected int    // kinds of objects
	}{
		{"single",
			"abc: x",
			1,
		},
		{"multiple",
			"abc: x\n---\ndef: y",
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			manifests, err := parseManifests(tc.input, "default")

			if err != nil {
				t.Fatalf("Expected to successfully parse manifests: %v", err)
			}

			if len(manifests) != tc.expected {
				t.Fatalf("Expected %d manifests, instead got: %d", tc.expected, len(manifests))
			}

		})
	}
}

func TestFixNamespace(t *testing.T) {
	testCases := []struct {
		name     string
		input    MetaOject // manifest
		expected string    // kinds of objects
	}{
		{"present",
			MetaOject{ObjectMeta: metav1.ObjectMeta{Namespace: "some-namespace"}},
			"some-namespace",
		},
		{"missing",
			MetaOject{},
			"default-namespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			fixNamespace(&tc.input, tc.expected)

			actual := tc.input.ObjectMeta.Namespace

			if actual != tc.expected {
				t.Fatalf("Expected namespace to be: %s, instead got: %s", tc.expected, actual)
			}

		})
	}
}
