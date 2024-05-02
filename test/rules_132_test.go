package test

import (
	"testing"
)

func TestRego132(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"FlowSchema", []string{"../fixtures/flowschema-v1beta3.yaml"}, []string{"FlowSchema"}},
		{"PriorityLevelConfiguration", []string{"../fixtures/prioritylevelconfiguration-v1beta3.yaml"}, []string{"PriorityLevelConfiguration"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
