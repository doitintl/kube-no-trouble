package test

import (
	"testing"
)

func TestRego126(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"AutoScaler", []string{"../fixtures/autoscaler-v2beta2.yaml"}, []string{"HorizontalPodAutoscaler"}},
		{"FlowSchema", []string{"../fixtures/flowschema-v1beta1.yaml"}, []string{"FlowSchema"}},
		{"PriorityLevelConfiguration", []string{"../fixtures/prioritylevelconfiguration-v1beta1.yaml"}, []string{"PriorityLevelConfiguration"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
