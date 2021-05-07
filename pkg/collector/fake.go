package collector

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

func (c *fakeCollector) GetServerVersion() (string, error) {
	return FAKE_VERSION, nil
}
