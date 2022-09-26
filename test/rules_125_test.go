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
		{"CronJob", []string{"../fixtures/cronjob-v1beta1.yaml"}, []string{"CronJob"}},
		{"HorizontalPodAutoscaler", []string{"../fixtures/hpa-v2beta1.yaml"}, []string{"HorizontalPodAutoscaler"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
