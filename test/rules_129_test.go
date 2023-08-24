package test

import (
	"testing"
)

func TestRego129(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"FlowSchema", []string{"../fixtures/flowschema-v1beta2.yaml"}, []string{"FlowSchema"}},
		{"PriorityLevelConfiguration", []string{"../fixtures/prioritylevelconfiguration-v1beta2.yaml"}, []string{"PriorityLevelConfiguration"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
