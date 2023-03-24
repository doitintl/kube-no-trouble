package printer

import (
	"testing"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

var testVersion1, _ = judge.NewVersion("1.1.1")
var testVersion2, _ = judge.NewVersion("2.2.2")

var testInput []judge.Result = []judge.Result{
	{
		Name:        "testName1",
		Kind:        "testKind1",
		Namespace:   "testNamespace1",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset1",
		ReplaceWith: "testReplaceWith1",
		Since:       testVersion1,
	},
	{
		Name:        "testName2",
		Kind:        "testKind2",
		Namespace:   "testNamespace2",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset2",
		ReplaceWith: "testReplaceWith2",
		Since:       testVersion2,
	},
	{
		Name:        "testName3",
		Kind:        "testKind3",
		Namespace:   "testNamespace3",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset3",
		ReplaceWith: "testReplaceWith3",
		Since:       nil,
	},
}

func TestFilterNonRelevantResults(t *testing.T) {
	filterVersion, _ := judge.NewVersion("2.0.0")

	results, err := FilterNonRelevantResults(testInput[0:2], filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result after filter, got %d instead", len(results))
	}
}

func TestFilterNonRelevantResultsEmpty(t *testing.T) {
	filterVersion, _ := judge.NewVersion("2.0.0")

	var input []judge.Result = []judge.Result{}

	results, err := FilterNonRelevantResults(input, filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results after filter, got %d instead", len(results))
	}
}

func TestFilterNonRelevantResultsWithNilVersion(t *testing.T) {
	filterVersion, _ := judge.NewVersion("2.0.0")

	results, err := FilterNonRelevantResults(testInput[2:3], filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 results after filter, got %d instead", len(results))
	}
}

func TestFilterNonRelevantResultsNilTargetVersion(t *testing.T) {
	var filterVersion *judge.Version

	results, err := FilterNonRelevantResults(testInput, filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != len(testInput) {
		t.Errorf("expected same number of results as in input: %d, got %d instead", len(testInput), len(results))
	}
}
