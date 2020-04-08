package collector

type Collector interface {
	Get() ([]interface{}, error)
	Name() string
}

type commonCollector struct {
	name string
}

func (c *commonCollector) Name() string {
	return c.name
}
