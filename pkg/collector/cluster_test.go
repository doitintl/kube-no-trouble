package collector

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	discoveryFake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestNewClusterCollectorBadPath(t *testing.T) {
	testOpts := ClusterOpts{Kubeconfig: "bad path"}
	_, err := NewClusterCollector(&testOpts, []string{}, []string{}, USER_AGENT)

	if !strings.Contains(err.Error(), "no configuration has been provided") {
		if err != nil {
			t.Errorf("Should have errored with invalid configuration error instead got: %s", err)
		} else {
			t.Errorf("Should have failed but succeeded")
		}
	}
}

func TestNewClusterCollectorValidEmptyCollector(t *testing.T) {
	scheme := runtime.NewScheme()
	clientset := fake.NewSimpleDynamicClient(scheme)
	discoveryClient := discoveryFake.FakeDiscovery{}
	testOpts := ClusterOpts{
		Kubeconfig:      filepath.Join(FIXTURES_DIR, "kube.config"),
		ClientSet:       clientset,
		DiscoveryClient: &discoveryClient,
	}
	collector, err := NewClusterCollector(&testOpts, []string{}, []string{}, USER_AGENT)

	if err != nil {
		t.Fatalf("Should have parsed config instead got: %s", err)
	}

	result, err := collector.Get()

	if err != nil && result != nil {
		t.Errorf("Invalid schema")
	}
}

func TestNewClusterCollectorFakeClient(t *testing.T) {
	scheme := runtime.NewScheme()
	clientset := fake.NewSimpleDynamicClient(scheme)
	discoveryClient := discoveryFake.FakeDiscovery{}
	testOpts := ClusterOpts{ClientSet: clientset, DiscoveryClient: &discoveryClient}

	collector, err := NewClusterCollector(&testOpts, []string{}, []string{}, USER_AGENT)
	if err != nil {
		t.Fatalf("failed to create cluster collector from fake client: %s", err)
	}

	result, err := collector.Get()

	if err != nil || len(result) != 0 {
		t.Errorf("expected to receive zero resources")
	}
}

func TestClusterCollectorGetFake(t *testing.T) {
	testCases := []struct {
		name                  string
		input                 []string // file list
		additionalAnnotations []string
		expected              int // number of manifests
	}{
		{
			name:     "empty",
			input:    []string{},
			expected: 0,
		},
		{
			name:     "withoutAnnotation",
			input:    []string{"fake-deployment-v1beta1-no-annotation.yaml"},
			expected: 1,
		},
		{
			name:     "one",
			input:    []string{"fake-deployment-v1beta1-with-annotation.yaml"},
			expected: 1,
		},
		{
			name:     "multiple",
			input:    []string{"fake-deployment-v1beta1-with-annotation.yaml", "fake-ingress-v1beta1-with-annotation.yaml"},
			expected: 2,
		},
		{
			name:     "mixed",
			input:    []string{"fake-deployment-v1beta1-no-annotation.yaml", "fake-ingress-v1beta1-with-annotation.yaml"},
			expected: 2,
		},
		{
			name:                  "kappAnnotation",
			input:                 []string{"fake-deployment-v1beta1-with-kapp-annotation.yaml"},
			additionalAnnotations: []string{"kapp.k14s.io/original"},
			expected:              2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//var objs []*unstructured.Unstructured
			var objs []runtime.Object
			for _, f := range tc.input {
				obj := &unstructured.Unstructured{}

				input, err := ioutil.ReadFile(filepath.Join(FIXTURES_DIR, f))
				if err != nil {
					t.Errorf("failed to read file: %s", err)
				}
				dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				_, _, err = dec.Decode(input, nil, obj)
				if err != nil {
					t.Errorf("failed to decode input manifest %s: %s", f, err)
				}
				objs = append(objs, obj)
			}

			rscheme := runtime.NewScheme()
			_ = scheme.AddToScheme(rscheme)

			clientset := fake.NewSimpleDynamicClient(rscheme, objs...)
			discoveryClient := discoveryFake.FakeDiscovery{}
			testOpts := ClusterOpts{ClientSet: clientset, DiscoveryClient: &discoveryClient}

			collector, err := NewClusterCollector(&testOpts, []string{}, tc.additionalAnnotations, USER_AGENT)

			if err != nil {
				t.Errorf("failed to create collector from fake client: %s", err)
			}

			result, err := collector.Get()

			if err != nil {
				t.Errorf("expected to receive resources: %s", err)
			}
			if len(result) != tc.expected {
				t.Errorf("expected to receive %d, received %d resources", tc.expected, len(result))
			}
		})
	}
}
