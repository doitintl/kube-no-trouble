package test

import (
	"testing"
)

func TestRego126(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"AutoScaler-v2beta1", []string{"../fixtures/autoscaler-v2beta1.yaml"}, []string{"HorizontalPodAutoscaler"}},
		{"AutoScaler-v2beta2", []string{"../fixtures/autoscaler-v2beta2.yaml"}, []string{"HorizontalPodAutoscaler"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
