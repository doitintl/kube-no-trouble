package test

import (
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
)

func TestRego122(t *testing.T) {
	testCases := []struct {
		name          string
		manifests     []string
		expectedKinds []string // kinds of objects
	}{
		{"ClusterRole", []string{"../fixtures/clusterrole-v1beta1.yaml"}, []string{"ClusterRole"}},
		{"ClusterRoleBinding", []string{"../fixtures/clusterrolebinding-v1beta1.yaml"}, []string{"ClusterRoleBinding"}},
		{"CSIDriver", []string{"../fixtures/csidriver-v1beta1.yaml"}, []string{"CSIDriver"}},
		{"CSINode", []string{"../fixtures/csinode-v1beta1.yaml"}, []string{"CSINode"}},
		{"Role", []string{"../fixtures/role-v1beta1.yaml"}, []string{"Role"}},
		{"RoleBinding", []string{"../fixtures/rolebinding-v1beta1.yaml"}, []string{"RoleBinding"}},
		{"StorageClass", []string{"../fixtures/storageclass-v1beta1.yaml"}, []string{"StorageClass"}},
		{"VolumeAttachment", []string{"../fixtures/volumeattachment-v1beta1.yaml"}, []string{"VolumeAttachment"}},
		{"PriorityClass", []string{"../fixtures/priorityclass-v1beta1.yaml"}, []string{"PriorityClass"}},
		{"Ingress", []string{"../fixtures/ingress-v1beta1.yaml"}, []string{"Ingress"}},
		{"IngressClass", []string{"../fixtures/ingressclass-v1beta1.yaml"}, []string{"IngressClass"}},
		{"Lease", []string{"../fixtures/lease-v1beta1.yaml"}, []string{"Lease"}},
		{"CertificateSigningRequest", []string{"../fixtures/certificatesigningrequest-v1beta1.yaml"}, []string{"CertificateSigningRequest"}},
		{"SubjectAccessReview", []string{"../fixtures/subjectaccessreview-v1beta1.yaml"}, []string{"SubjectAccessReview"}},
		{"SelfSubjectAccessReview", []string{"../fixtures/selfsubjectaccessreview-v1beta1.yaml"}, []string{"SelfSubjectAccessReview"}},
		{"LocalSubjectAccessReview", []string{"../fixtures/localsubjectaccessreview-v1beta1.yaml"}, []string{"LocalSubjectAccessReview"}},
		{"TokenReview", []string{"../fixtures/tokenreview-v1beta1.yaml"}, []string{"TokenReview"}},
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
