package collector

import (
	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

const (
	FAKE_VERSION        = "1.2.3"
	FAKE_COLLECTOR_NAME = "Fake"
)

type fakeCollector struct {
	*commonCollector
}

func (c *fakeCollector) Get() ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func NewFakeCollector() *fakeCollector {
	return &fakeCollector{
		commonCollector: newCommonCollector(FAKE_COLLECTOR_NAME),
	}
}

func (c *fakeCollector) GetServerVersion() (*judge.Version, error) {
	version, err := judge.NewVersion(FAKE_VERSION)

	return version, err
}
