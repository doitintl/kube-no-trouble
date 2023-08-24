package test

import (
	"testing"
)

func TestRego127(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"CSIStorageCapacity", []string{"../fixtures/csistoragecapacity-v1beta1.yaml"}, []string{"CSIStorageCapacity"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
