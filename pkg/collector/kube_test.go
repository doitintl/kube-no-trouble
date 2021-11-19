package collector

import (
	"k8s.io/client-go/rest"
	"os"
	"path/filepath"
	"testing"

	goversion "github.com/hashicorp/go-version"
	"k8s.io/apimachinery/pkg/version"
	discoveryFake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
)

const FIXTURES_DIR = "../../fixtures"

func TestNewKubeCollector(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	col, err := newKubeCollector("", "", clientSet.Discovery())

	if err != nil {
		t.Fatalf("Failed to create kubeCollector from fake discovery client")
	}
	if col == nil {
		t.Fatalf("Should return collector, got nil instead")
	}
}

func TestNewKubeCollectorWithKubeconfigPath(t *testing.T) {
	col, err := newKubeCollector(filepath.Join(FIXTURES_DIR, "kube.config.basic"), "", nil)

	if err != nil {
		t.Fatalf("Failed to create kubeCollector from fake discovery client")
	}
	if col == nil {
		t.Fatalf("Should return collector, got nil instead")
	}
	if col.discoveryClient == nil {
		t.Fatalf("Collector should have discovery client, got nil instead")
	}
	if col.restConfig == nil {
		t.Fatalf("Collector should have rest config, got nil instead")
	}
}

func TestNewKubeCollectorError(t *testing.T) {
	_, err := newKubeCollector("does-not-exist", "", nil)

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

	collector, err := newKubeCollector("", "", clientSet.Discovery())
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

func TestContext(t *testing.T) {
	expectedContext := "test-context"
	expectedHost := "https://test-cluster"

	collector, err := newKubeCollector(filepath.Join(FIXTURES_DIR, "kube.config.context"), expectedContext, nil)
	if err != nil {
		t.Fatalf("Failed to create kubeCollector from fake client with context %s: %s", expectedContext, err)
	}

	host := collector.GetRestConfig().Host
	if host != expectedHost {
		t.Fatalf("Expected host from context %s to be: %s, got %s instead", expectedContext, expectedHost, host)
	}
}

func TestContextMissing(t *testing.T) {
	expectedContext := "non-existent"

	_, err := newKubeCollector(filepath.Join(FIXTURES_DIR, "kube.config.context"), expectedContext, nil)
	if err == nil {
		t.Fatalf("Expected to fail when uisng non-existent context: %s", expectedContext)
	}
}

func TestNewClientRestConfigError(t *testing.T) {
	_, err := newClientRestConfig("does-not-exist", "", rest.InClusterConfig)

	if err == nil {
		t.Errorf("Expected to fail with non-existent kubeconfig")
	}
}

func TestNewClientRestConfigWithKubeconfigPath(t *testing.T) {
	_, err := newClientRestConfig(filepath.Join(FIXTURES_DIR, "kube.config.basic"), "", rest.InClusterConfig)
	if err != nil {
		t.Fatalf("Failed with: %s", err)
	}
}

func TestNewClientRestConfigMultipleFiles(t *testing.T) {
	env := setupEnv(t, map[string]string{
		"KUBECONFIG": filepath.Join(FIXTURES_DIR, "kube.config.empty") + ":" + filepath.Join(FIXTURES_DIR, "kube.config.basic"),
	})
	defer env.reset()

	if _, err := newClientRestConfig("", "", rest.InClusterConfig); err != nil {
		t.Fatalf("Failed with: %s", err)
	}
}

func TestNewClientRestConfigWithContext(t *testing.T) {
	expectedContext := "test-context"
	expectedHost := "https://test-cluster"

	config, err := newClientRestConfig(filepath.Join(FIXTURES_DIR, "kube.config.context"), expectedContext, rest.InClusterConfig)
	if err != nil {
		t.Fatalf("Failed to create kubeCollector from fake client with context %s: %s", expectedContext, err)
	}

	if config.Host != expectedHost {
		t.Fatalf("Expected host from context %s to be: %s, got %s instead", expectedContext, expectedHost, config.Host)
	}
}

func TestNewClientRestConfigContextMissing(t *testing.T) {
	expectedContext := "non-existent"

	_, err := newClientRestConfig(filepath.Join(FIXTURES_DIR, "kube.config.context"), expectedContext, rest.InClusterConfig)
	if err == nil {
		t.Fatalf("Expected to fail when uisng non-existent context: %s", expectedContext)
	}
}

func TestNewClientRestConfigInCluster(t *testing.T) {
	expectedHost := "test-host"
	inClusterFn := func() (*rest.Config, error) {
		return &rest.Config{Host: expectedHost}, nil
	}

	cfg, err := newClientRestConfig("", "", inClusterFn)
	if err != nil {
		t.Fatalf("Failed to create in-cluster config: %s", err)
	}
	if cfg.Host != expectedHost {
		t.Fatalf("Expected %s host, instead got: %s", expectedHost, cfg.Host)
	}
}

type env struct {
	initialVals  map[string]string
	initialState map[string]bool
	t            *testing.T
}

func setupEnv(t *testing.T, vars map[string]string) *env {
	env := env{
		initialState: make(map[string]bool),
		initialVals:  make(map[string]string),
	}

	for k, v := range vars {
		env.initialVals[k], env.initialState[k] = os.LookupEnv(k)
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("Failed so set %s env variable for test: %s", k, err)
		}
	}

	return &env
}

func (e *env) reset() {
	for k, v := range e.initialState {
		var err error
		if v {
			err = os.Setenv(k, e.initialVals[k])
		} else {
			err = os.Unsetenv(k)
		}

		if err != nil {
			e.t.Errorf("Failed to reset %s env variable after test: %s", k, err)
		}
	}
}
