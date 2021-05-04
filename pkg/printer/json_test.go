package printer

import (
	"encoding/json"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

var findingsJsonTesting []judge.Result = []judge.Result{
	{
		Name:        "testName1",
		Kind:        "testKind1",
		Namespace:   "testNamespace1",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset1",
		ReplaceWith: "testReplaceWith1",
		Since:       "testSince1",
	},
	{
		Name:        "testName2",
		Kind:        "testKind2",
		Namespace:   "testNamespace2",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset2",
		ReplaceWith: "testReplaceWith2",
		Since:       "testSince2",
	},
}

func TestJsonPopulateOutput(t *testing.T) {
	printer := &jsonPrinter{}
	output, err := printer.populateOutput(findingsJsonTesting)

	if err != nil {
		t.Fatal(err)
	}

	var j []judge.Result
	err = json.Unmarshal([]byte(output), &j)
	if err != nil {
		t.Fatal(err)
	}
	if len(j) != 2 {
		t.Error("wrong number of results")
	}
}
