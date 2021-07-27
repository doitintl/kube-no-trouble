package collector

import (
	"reflect"
	"testing"
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
		input    map[string]interface{} // manifest
		expected string                 // kinds of objects
	}{
		{"present",
			map[string]interface{}{"metadata": map[string]interface{}{"namespace": "some-namespace"}},
			"some-namespace",
		},
		{"missing",
			map[string]interface{}{"metadata": map[string]interface{}{}},
			"default-namespace",
		},
		{"nil",
			map[string]interface{}{"metadata": map[string]interface{}{"namespace": nil}},
			"default-namespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			fixNamespace(&tc.input, tc.expected)

			meta, ok := tc.input["metadata"]
			if !ok {
				t.Fatalf("Expected fixed manifest to have metadata key")
			}

			expectedType := "map[string]interface {}"
			if reflect.TypeOf(meta).String() != expectedType {
				t.Fatalf("Expected metadata type to be %s, instead got: %T", expectedType, meta)
			}

			actual, ok := meta.(map[string]interface{})["namespace"]
			if !ok {
				t.Fatalf("Expected fixed manifest to have metadata.namespace key")
			}

			if actual != tc.expected {
				t.Fatalf("Expected namespace to be: %s, instead got: %s", tc.expected, actual)
			}

		})
	}
}
