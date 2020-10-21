package judge

import (
	"testing"
)

func TestEvalEmpty(t *testing.T) {
	inputs := []map[string]interface{}{}

	judge, err := NewRegoJudge(&RegoOpts{})
	if err != nil {
		t.Errorf("failed to create judge instance: %s", err)
	}

	results, err := judge.Eval(inputs)
	if err != nil {
		t.Errorf("failed to evaluate input: %s", err)
	}

	if results == nil || len(results) != 0 {
		t.Errorf("expected empty array, instead got: %v", results)
	}
}
