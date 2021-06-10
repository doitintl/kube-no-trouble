package test

import (
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
)

func TestRego116(t *testing.T) {
	testCases := []struct {
		name          string
		manifests     []string
		expectedKinds []string // kinds of objects
	}{
		{"CSIDriver", []string{"../fixtures/csidriver-v1beta1.yaml"}, []string{"CSIDriver"}},
		{"CSINode", []string{"../fixtures/csinode-v1beta1.yaml"}, []string{"CSINode"}},
		{"StorageClass", []string{"../fixtures/storageclass-v1beta1.yaml"}, []string{"StorageClass"}},
		{"VolumeAttachment", []string{"../fixtures/volumeattachment-v1beta1.yaml"}, []string{"VolumeAttachment"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := collector.NewFileCollector(
				&collector.FileOpts{Filenames: tc.manifests},
			)

			if err != nil {
				t.Errorf("Expected to succeed for %s, failed: %s", tc.manifests, err)
			}

			manifests, err := c.Get()
			if err != nil {
				t.Errorf("Expected to succeed for %s, failed: %s", tc.manifests, err)
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
