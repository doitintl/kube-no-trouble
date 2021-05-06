package collector

import (
	"testing"

	"k8s.io/apimachinery/pkg/version"
	discoveryFake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewKubeCollector(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	col, err := newKubeCollector("", clientSet.Discovery())

	if err != nil {
		t.Errorf("Failed to create kubeCollector form fake discovery client")
	}
	if col == nil {
		t.Errorf("Should return collector, instrad got nil")
	}

}

func TestNewKubeCollectorError(t *testing.T) {
	_, err := newKubeCollector("does-not-exist", nil)

	if err == nil {
		t.Errorf("Expected to fail with non-existent kubeconfig")
	}
}

func TestGetServerVersion(t *testing.T) {
	expectedMajor := "1"
	expectedMinor := "2"
	expectedVersion := expectedMajor + "." + expectedMinor

	clientSet := fake.NewSimpleClientset()
	clientSet.Discovery().(*discoveryFake.FakeDiscovery).FakedServerVersion = &version.Info{
		Major: expectedMajor,
		Minor: expectedMinor,
	}

	collector, err := newKubeCollector("", clientSet.Discovery())
	if err != nil {
		t.Fatalf("failed to create kubeCollector from fake client: %s", err)
	}

	version, err := collector.GetServerVersion()
	if err != nil {
		t.Errorf("Failed to get version with error: %s", err)
	}

	if version != expectedVersion {
		t.Errorf("Expected no version to be detected, instead got: %s", version)
	}
}
