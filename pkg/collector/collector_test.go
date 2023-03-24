package collector

import (
	"testing"

	"github.com/LeMyst/kube-no-trouble/test"
)

func init() {
	test.Setup()
}

func TestName(t *testing.T) {
	testCollector := commonCollector{name: "I am a collector"}
	result := testCollector.Name()

	if result != "I am a collector" {
		t.Errorf("Collector name required")
	}
}
