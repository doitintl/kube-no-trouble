package collector

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	discoveryFake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestNewClusterCollectorBadPath(t *testing.T) {
	testOpts := ClusterOpts{Kubeconfig: "bad path"}
	_, err := NewClusterCollector(&testOpts, []string{}, []string{}, []string{}, USER_AGENT)

	if !strings.Contains(err.Error(), "no configuration has been provided") {
		if err != nil {
			t.Errorf("Should have errored with invalid configuration error instead got: %s", err)
		} else {
			t.Errorf("Should have failed but succeeded")
		}
	}
}

func TestNewClusterCollectorValidEmptyCollector(t *testing.T) {
	rscheme := scheme.Scheme
	registerRequiredListsForFakeClient(rscheme)

	clientset := fake.NewSimpleDynamicClient(rscheme)

	discoveryClient := discoveryFake.FakeDiscovery{}
	testOpts := ClusterOpts{
		Kubeconfig:      filepath.Join(FIXTURES_DIR, "kube.config"),
		ClientSet:       clientset,
		DiscoveryClient: &discoveryClient,
	}
	collector, err := NewClusterCollector(&testOpts, []string{}, []string{}, []string{}, USER_AGENT)

	if err != nil {
		t.Fatalf("Should have parsed config instead got: %s", err)
	}

	result, err := collector.Get()

	if err != nil && result != nil {
		t.Errorf("Invalid schema")
	}
}

func TestNewClusterCollectorFakeClient(t *testing.T) {
	rscheme := scheme.Scheme
	registerRequiredListsForFakeClient(rscheme)

	clientset := fake.NewSimpleDynamicClientWithCustomListKinds(rscheme, map[schema.GroupVersionResource]string{})
	discoveryClient := discoveryFake.FakeDiscovery{}
	testOpts := ClusterOpts{ClientSet: clientset, DiscoveryClient: &discoveryClient}

	collector, err := NewClusterCollector(&testOpts, []string{}, []string{}, []string{}, USER_AGENT)
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
			expected: 0,
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
			expected: 1,
		},
		{
			name:                  "kappAnnotation",
			input:                 []string{"fake-deployment-v1beta1-with-kapp-annotation.yaml"},
			additionalAnnotations: []string{"kapp.k14s.io/original"},
			expected:              1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//var objs []*unstructured.Unstructured
			var objs []runtime.Object
			for _, f := range tc.input {
				obj := &unstructured.Unstructured{}

				input, err := ioutil.ReadFile(filepath.Join(FIXTURES_DIR, f))
				dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				_, _, err = dec.Decode(input, nil, obj)
				if err != nil {
					t.Errorf("failed to decode input manifest %s: %s", f, err)
				}
				objs = append(objs, obj)
			}

			rscheme := scheme.Scheme
			registerRequiredListsForFakeClient(rscheme)

			clientset := fake.NewSimpleDynamicClient(rscheme, objs...)
			discoveryClient := discoveryFake.FakeDiscovery{}
			testOpts := ClusterOpts{ClientSet: clientset, DiscoveryClient: &discoveryClient}

			namespaces := []string{""}
			collector, err := NewClusterCollector(&testOpts, namespaces, []string{}, tc.additionalAnnotations, USER_AGENT)

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

func TestClusterCollector_getLastAppliedConfig(t *testing.T) {
	tests := []struct {
		name                string
		c                   *ClusterCollector
		resourceAnnotations map[string]string
		wantManifest        string
		wantOk              bool
	}{
		{
			name: "Default annotation",
			c:    &ClusterCollector{},
			resourceAnnotations: map[string]string{
				"kubectl.kubernetes.io/last-applied-configuration": "Some config",
				"another-annotation": "Bla",
			},
			wantManifest: "Some config",
			wantOk:       true,
		},
		{
			name: "Try kubectl annotation first",
			c: &ClusterCollector{
				additionalAnnotations: []string{"kapp.k14s.io/original"},
			},
			resourceAnnotations: map[string]string{
				"kubectl.kubernetes.io/last-applied-configuration": "Some config",
				"kapp.k14s.io/original":                            "Kapp config",
				"another-annotation":                               "Bla",
			},
			wantManifest: "Some config",
			wantOk:       true,
		},
		{
			name: "Use additional annotation",
			c: &ClusterCollector{
				additionalAnnotations: []string{"kapp.k14s.io/original"},
			},
			resourceAnnotations: map[string]string{
				"kapp.k14s.io/original": "Kapp config",
				"another-annotation":    "Bla",
			},
			wantManifest: "Kapp config",
			wantOk:       true,
		},
		{
			name: "No annotation found",
			c: &ClusterCollector{
				additionalAnnotations: []string{},
			},
			resourceAnnotations: map[string]string{
				"kapp.k14s.io/original": "Kapp config",
				"another-annotation":    "Bla",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest, ok := tt.c.getLastAppliedConfig(tt.resourceAnnotations)
			if manifest != tt.wantManifest {
				t.Errorf("ClusterCollector.getLastAppliedConfig() got = %v, want %v", manifest, tt.wantManifest)
			}
			if ok != tt.wantOk {
				t.Errorf("ClusterCollector.getLastAppliedConfig() got1 = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func registerRequiredListsForFakeClient(s *runtime.Scheme) {
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "authorization.k8s.io", Version: "v1", Kind: "SubjectAccessReviewList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "authorization.k8s.io", Version: "v1", Kind: "SelfSubjectAccessReviewList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "authorization.k8s.io", Version: "v1", Kind: "LocalSubjectAccessReviewList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "authentication.k8s.io", Version: "v1", Kind: "TokenReviewList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "apiregistration.k8s.io", Version: "v1", Kind: "ApiServiceList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "snapshot.storage.k8s.io", Version: "v1", Kind: "VolumeSnapshotList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "snapshot.storage.k8s.io", Version: "v1", Kind: "VolumeSnapshotClassList"}, &unstructured.UnstructuredList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "snapshot.storage.k8s.io", Version: "v1", Kind: "VolumeSnapshotContentList"}, &unstructured.UnstructuredList{})
}
