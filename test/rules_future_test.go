package test

import (
	"testing"
)

func TestRegoFuture(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"VolumeSnapshot", []string{"../fixtures/volumesnapshot-v1beta1.yaml"}, []string{"VolumeSnapshot"}},
		{"VolumeSnapshotClass", []string{"../fixtures/volumesnapshotclass-v1beta1.yaml"}, []string{"VolumeSnapshotClass"}},
		{"VolumeSnapshotContent", []string{"../fixtures/volumesnapshotcontent-v1beta1.yaml"}, []string{"VolumeSnapshotContent"}},
	}

	testReourcesUsingFixtures(t, testCases)
}
