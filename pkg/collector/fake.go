package collector

const (
	FAKE_VERSION        = "1.2.3"
	FAKE_COLLECTOR_NAME = "Fake"
)

type fakeCollector struct {
	*commonCollector
}

func (c *fakeCollector) Get() ([]MetaOject, error) {
	return []MetaOject{}, nil
}

func NewFakeCollector() *fakeCollector {
	return &fakeCollector{
		commonCollector: newCommonCollector(FAKE_COLLECTOR_NAME),
	}
}

func (c *fakeCollector) GetServerVersion() (*Version, error) {
	version, err := NewVersion(FAKE_VERSION)

	return version, err
}
