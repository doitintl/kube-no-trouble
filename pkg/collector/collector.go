package collector

type Collector interface {
	Get() ([]map[string]interface{}, error)
	Name() string
}

type commonCollector struct {
	name string
}

func (c *commonCollector) Name() string {
	return c.name
}
