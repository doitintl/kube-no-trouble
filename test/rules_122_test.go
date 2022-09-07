package test

import (
	"testing"
)

func TestRego122(t *testing.T) {
	testCases := []resourceFixtureTestCase{
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
		{"APIService", []string{"../fixtures/apiservice-v1beta1.yaml"}, []string{"APIService"}},
		{"CustomResourceDefinition", []string{"../fixtures/customresourcedefinition-v1beta1.yaml"}, []string{"CustomResourceDefinition"}},
		{"MutatingWebhookConfiguration", []string{"../fixtures/mutatingwebhookconfiguration-v1beta1.yaml"}, []string{"MutatingWebhookConfiguration"}},
		{"ValidatingWebhookConfiguration", []string{"../fixtures/validatingwebhookconfiguration-v1beta1.yaml"}, []string{"ValidatingWebhookConfiguration"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
