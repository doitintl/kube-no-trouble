package main

import (
	"errors"
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
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	initCollectors := []collector.Collector{}
	initCollectors = append(initCollectors, fileCollector)

	collectors := getCollectors(initCollectors)

	if collectors != nil && len(collectors) != 1 {
		t.Errorf("Did not get file collector correctly with error: %s", err)
	}
}

func TestStoreCollector(t *testing.T) {
	collectors := []collector.Collector{}
	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})

	if err != nil {
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	collectors = storeCollector(fileCollector, err, collectors)

	if len(collectors) != 1 {
		t.Errorf("Failed to append collector")
	}
}

func TestStoreCollectorMultiple(t *testing.T) {
	collectors := []collector.Collector{}
	fileCollector, err := collector.NewFileCollector(&collector.FileOpts{Filenames: []string{"../../fixtures/deployment-v1beta1.yaml"}})

	if err != nil {
		t.Errorf("Failed to create File collector with error: %s", err)
	}

	collectors = storeCollector(fileCollector, err, collectors)

	collectors = storeCollector(fileCollector, err, collectors)

	if len(collectors) != 2 {
		t.Errorf("Failed to append collectors")
	}
}

func TestStoreCollectorError(t *testing.T) {
	collectors := []collector.Collector{}
	err := errors.New("Just testing...")

	collectors = storeCollector(nil, err, collectors)

	if len(collectors) != 0 {
		t.Errorf("Failed to ignore collector with error")
	}
}
