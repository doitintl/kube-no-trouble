package collector

import (
	goversion "github.com/hashicorp/go-version"
)

type Collector interface {
	Get() ([]map[string]interface{}, error)
	Name() string
}

type VersionCollector interface {
	GetServerVersion() (*goversion.Version, error)
}

type commonCollector struct {
	name string
}

func newCommonCollector(name string) *commonCollector {
	return &commonCollector{
		name: name,
	}
}

func (c *commonCollector) Name() string {
	return c.name
}
