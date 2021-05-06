package collector

type Collector interface {
	Get() ([]map[string]interface{}, error)
	Name() string
}

type VersionCollector interface {
	GetServerVersion() (string, error)
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
