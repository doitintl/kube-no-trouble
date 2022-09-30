package collector

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

var FAKE_API_RESOURCES = &metav1.APIResourceList{
	GroupVersion: "v1",
	APIResources: []metav1.APIResource{
		{Name: "pods", Namespaced: true, Kind: "Pod"},
		{Name: "services", Namespaced: true, Kind: "Service"},
		{Name: "namespaces", Namespaced: false, Kind: "Namespace"},
	},
}

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

			fcs := fake.NewSimpleClientset()
			manifests, err := parseManifests(tc.input, "default", fcs.Discovery())

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
			map[string]interface{}{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]interface{}{"namespace": "some-namespace"}},
			"some-namespace",
		},
		{"missing",
			map[string]interface{}{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]interface{}{}},
			"default-namespace",
		},
		{"nil",
			map[string]interface{}{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]interface{}{"namespace": nil}},
			"default-namespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			fcs := fake.NewSimpleClientset()
			fcs.Resources = []*metav1.APIResourceList{FAKE_API_RESOURCES}
			fixNamespace(&unstructured.Unstructured{tc.input}, tc.expected, fcs.Discovery())

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

func Test_isResourceNamespaced(t *testing.T) {
	tests := []struct {
		name string
		gvk  string
		want bool
	}{
		{"true-pod", "Pod.v1.", true},
		{"false-namespace", "Namespace.v1.", false},
		{"false-non-existent", "XXX.v1.", false},
	}

	fcs := fake.NewSimpleClientset()
	fcs.Resources = []*metav1.APIResourceList{FAKE_API_RESOURCES}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gvk, _ := schema.ParseKindArg(tt.gvk)

			got := isResourceNamespaced(fcs.Discovery(), *gvk)
			if got != tt.want {
				t.Errorf("isResourceNamespaced() got = %v, want %v", got, tt.want)
			}
		})
	}
}
