package test

import (
	"testing"
)

func TestRegoFuture(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"VolumeSnapshot", []string{"../fixtures/volumesnapshot-v1beta1.yaml"}, []string{"VolumeSnapshot"}},
		{"VolumeSnapshotClass", []string{"../fixtures/volumesnapshotclass-v1beta1.yaml"}, []string{"VolumeSnapshotClass"}},
		{"VolumeSnapshotContent", []string{"../fixtures/volumesnapshotcontent-v1beta1.yaml"}, []string{"VolumeSnapshotContent"}},
		{"SelfSubjectReview", []string{"../fixtures/selfsubjectreview-v1beta1.yaml"}, []string{"SelfSubjectReview"}},
		{"ValidatingAdmissionPolicy", []string{"../fixtures/validatingadmissionpolicy-v1beta1.yaml"}, []string{"ValidatingAdmissionPolicy"}},
		{"ValidatingAdmissionPolicyBinding", []string{"../fixtures/validatingadmissionpolicybinding-v1beta1.yaml"}, []string{"ValidatingAdmissionPolicyBinding"}},
		{"IPAddress", []string{"../fixtures/ipaddress-v1beta1.yaml"}, []string{"IPAddress"}},
		{"ServiceCIDR", []string{"../fixtures/servicecidr-v1beta1.yaml"}, []string{"ServiceCIDR"}},
		{"VolumeAttributesClass", []string{"../fixtures/volumeattributesclass-v1beta1.yaml"}, []string{"VolumeAttributesClass"}},
		{"DeviceClass-v1beta1", []string{"../fixtures/deviceclass-v1beta1.yaml"}, []string{"DeviceClass"}},
		{"DeviceClass-v1beta2", []string{"../fixtures/deviceclass-v1beta2.yaml"}, []string{"DeviceClass"}},
		{"ResourceClaim-v1beta1", []string{"../fixtures/resourceclaim-v1beta1.yaml"}, []string{"ResourceClaim"}},
		{"ResourceClaim-v1beta2", []string{"../fixtures/resourceclaim-v1beta2.yaml"}, []string{"ResourceClaim"}},
		{"ResourceClaimTemplate-v1beta1", []string{"../fixtures/resourceclaimtemplate-v1beta1.yaml"}, []string{"ResourceClaimTemplate"}},
		{"ResourceClaimTemplate-v1beta2", []string{"../fixtures/resourceclaimtemplate-v1beta2.yaml"}, []string{"ResourceClaimTemplate"}},
		{"ResourceSlice-v1beta1", []string{"../fixtures/resourceslice-v1beta1.yaml"}, []string{"ResourceSlice"}},
		{"ResourceSlice-v1beta2", []string{"../fixtures/resourceslice-v1beta2.yaml"}, []string{"ResourceSlice"}},
		{"MutatingAdmissionPolicy", []string{"../fixtures/mutatingadmissionpolicy-v1beta1.yaml"}, []string{"MutatingAdmissionPolicy"}},
		{"MutatingAdmissionPolicyBinding", []string{"../fixtures/mutatingadmissionpolicybinding-v1beta1.yaml"}, []string{"MutatingAdmissionPolicyBinding"}},
		{"ClusterTrustBundle", []string{"../fixtures/clustertrustbundle-v1beta1.yaml"}, []string{"ClusterTrustBundle"}},
		{"PodCertificateRequest", []string{"../fixtures/podcertificaterequest-v1beta1.yaml"}, []string{"PodCertificateRequest"}},
		{"StorageVersionMigration", []string{"../fixtures/storageversionmigration-v1beta1.yaml"}, []string{"StorageVersionMigration"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
