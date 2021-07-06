package printer

import (
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/judge"

	goversion "github.com/hashicorp/go-version"
)

var testVersion1, _ = goversion.NewVersion("1.1.1")
var testVersion2, _ = goversion.NewVersion("2.2.2")

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
}

func TestFilterNonRelevantResults(t *testing.T) {
	filterVersion, _ := goversion.NewVersion("2.0.0")

	results, err := FilterNonRelevantResults(testInput, filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result after filter, got %d intead", len(results))
	}
}

func TestFilterNonRelevantResultsEmpty(t *testing.T) {
	filterVersion, _ := goversion.NewVersion("2.0.0")

	var input []judge.Result = []judge.Result{}

	results, err := FilterNonRelevantResults(input, filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results after filter, got %d intead", len(results))
	}
}

func TestFilterNonRelevantResultsNilVersion(t *testing.T) {
	var filterVersion *goversion.Version

	results, err := FilterNonRelevantResults(testInput, filterVersion)
	if err != nil {
		t.Fatalf("failed to filter results: %s", err)
	}

	if len(results) != len(testInput) {
		t.Errorf("expected same number of results as in input: %d, got %d intead", len(testInput), len(results))
	}
}
