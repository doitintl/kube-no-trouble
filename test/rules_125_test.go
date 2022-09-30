package test

import (
	"testing"
)

func TestRego125(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"CronJob", []string{"../fixtures/cronjob-v1beta1.yaml"}, []string{"CronJob"}},
		{"EndpointSlice", []string{"../fixtures/endpointslice-v1beta1.yaml"}, []string{"EndpointSlice"}},
		{"PodDisruptionBudget", []string{"../fixtures/poddisruptionbudget-v1beta1.yaml"}, []string{"PodDisruptionBudget"}},
		{"PodSecurityPolicy", []string{"../fixtures/podsecuritypolicy-v1beta1.yaml"}, []string{"PodSecurityPolicy"}},
		{"RuntimeClass", []string{"../fixtures/runtimeclass-v1beta1.yaml"}, []string{"RuntimeClass"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
