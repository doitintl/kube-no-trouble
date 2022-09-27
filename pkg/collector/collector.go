package collector

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetaOject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type Collector interface {
	Get() ([]MetaOject, error)
	Name() string
}

type VersionCollector interface {
	GetServerVersion() (*Version, error)
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
