package collector

import (
	"encoding/json"
	"testing"
)

func TestNewClusterCollectorBadPath(t *testing.T) {
	testOpts := ClusterOpts{Kubeconfig: "bad path"}
	result, funcErr := NewClusterCollector(&testOpts)

	if funcErr.Error() != "stat bad path: no such file or directory" {
		out, err := json.Marshal(result)
		if err != nil {
			t.Errorf("Should have crashed with path error instead got: %s", string(out))
		} else {
			t.Errorf("Should have crashed instead got un-parseable error: %s", funcErr)
		}
	}
}
