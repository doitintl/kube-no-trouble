package test

import (
	"testing"
)

func TestRego125(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"RuntimeClass", []string{"../fixtures/runtimeclass-v1beta1.yaml"}, []string{"RuntimeClass"}},
	}

	testReourcesUsingFixtures(t, testCases)
}
