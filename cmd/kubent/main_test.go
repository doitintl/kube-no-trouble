package main

import (
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/collector"
	"github.com/doitintl/kube-no-trouble/pkg/config"
)

func TestInitCollectors(t *testing.T) {
	testConfig := config.Config{
		Filenames:  []string{"../../fixtures/deployment-v1beta1.yaml"},
		Cluster:    false,
		Debug:      false,
		Helm2:      false,
		Helm3:      false,
		Kubeconfig: "test",
		Output:     "test",
	}

	collectors := initCollectors(&testConfig)

	if collectors[0].Name() != "File" {
		t.Errorf("Did not parse fixture with path %s", testConfig.Filenames[0])
	}
}

func TestGetCollectors(t *testing.T) {
	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})

	if err != nil {
		t.Errorf("Did not parse fixture %s, with error: %s", fileCollector.Name(), err)
	}

	initCollectors := []collector.Collector{}
	initCollectors = append(initCollectors, fileCollector)

	collectors := getCollectors(initCollectors)

	if collectors != nil && len(collectors) != 1 {
		t.Errorf("Did not get file collector correctly with error: %s", err)
	}
}
