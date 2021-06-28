package collector

import (
	goversion "github.com/hashicorp/go-version"
)

const FAKE_VERSION = "1.2.3"

type fakeCollector struct {
	*commonCollector
}

func (c *fakeCollector) Get() ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func NewFakeCollector() *fakeCollector {
	return &fakeCollector{
		commonCollector: newCommonCollector("Fake"),
	}
}

func (c *fakeCollector) GetServerVersion() (*goversion.Version, error) {
	version, err := goversion.NewVersion(FAKE_VERSION)

	return version, err
}
