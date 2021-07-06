package collector

import (
	"testing"

	goversion "github.com/hashicorp/go-version"
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
	gitVersion := "v1.2.3"
	expectedVersion, _ := goversion.NewVersion(gitVersion)

	clientSet := fake.NewSimpleClientset()
	clientSet.Discovery().(*discoveryFake.FakeDiscovery).FakedServerVersion = &version.Info{
		GitVersion: gitVersion,
	}

	collector, err := newKubeCollector("", clientSet.Discovery())
	if err != nil {
		t.Fatalf("failed to create kubeCollector from fake client: %s", err)
	}

	version, err := collector.GetServerVersion()
	if err != nil {
		t.Fatalf("Failed to get version with error: %s", err)
	}

	if !version.Equal(expectedVersion) {
		t.Errorf("Expected version: %s, instead got: %s", expectedVersion.String(), version.String())
	}
}
