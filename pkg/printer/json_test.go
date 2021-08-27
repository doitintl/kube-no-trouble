package printer

import (
	"encoding/json"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

func TestJsonPopulateOutput(t *testing.T) {
	printer := &jsonPrinter{}
	output, err := printer.populateOutput(testInput)

	if err != nil {
		t.Fatal(err)
	}

	var j []judge.Result
	err = json.Unmarshal([]byte(output), &j)
	if err != nil {
		t.Fatal(err)
	}
	if len(j) != len(testInput) {
		t.Error("wrong number of results")
	}
}
