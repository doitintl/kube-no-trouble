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
	if "1.1.1" != j[0].SinceStr {
		t.Errorf("Expected 1st Since '1.1.1' not '%s'", j[0].SinceStr)
	}
	if "1.1.2" != j[1].SinceStr {
		t.Errorf("Expected 2nd Since '1.1.2' not '%s'", j[1].SinceStr)
	}
	if "" != j[2].SinceStr {
		t.Errorf("Expected 3rd Since '' not '%s'", j[2].SinceStr)
	}
}
