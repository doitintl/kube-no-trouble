package test

import (
	"testing"
)

func TestRego126(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"AutoScaler", []string{"../fixtures/autoscaler-v2beta2.yaml"}, []string{"HorizontalPodAutoscaler"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
