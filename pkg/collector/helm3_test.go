package collector

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestNewHelmV3Collector(t *testing.T) {
	expectedName := "Helm v3"

	clientSet := fake.NewSimpleClientset()
	col, err := NewHelmV3Collector(&HelmV3Opts{
		DiscoveryClient: clientSet.Discovery(),
		CoreClient:      clientSet.CoreV1(),
	}, []string{}, USER_AGENT)

	if err != nil {
		t.Fatalf("Failed to create collector from fake discovery client")
	}
	if col == nil {
		t.Fatalf("Should return collector, got nil instead")
	}
	if col.Name() != expectedName {
		t.Fatalf("Expected collector name: %s, instead got: %s", expectedName, col.Name())
	}
}
