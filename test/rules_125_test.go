package test

import (
	"testing"
)

func TestRego125(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"RuntimeClass", []string{"../fixtures/runtimeclass-v1beta1.yaml"}, []string{"RuntimeClass"}},
		{"PodDisruptionBudget", []string{"../fixtures/poddisruptionbudget-v1beta1.yaml"}, []string{"PodDisruptionBudget"}},
		{"PodSecurityPolicy", []string{"../fixtures/podsecuritypolicy-v1beta1.yaml"}, []string{"PodSecurityPolicy"}},
		{"EndpointSlice", []string{"../fixtures/endpointslice-v1beta1.yaml"}, []string{"EndpointSlice"}},
	}

	testReourcesUsingFixtures(t, testCases)
}
