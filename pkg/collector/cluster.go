package collector

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterCollector struct {
	*commonCollector
	clientset dynamic.Interface
}

type ClusterOpts struct {
	Kubeconfig string
}

func NewClusterCollector(opts *ClusterOpts) (*ClusterCollector, error) {
	collector := &ClusterCollector{commonCollector: &commonCollector{name: "Cluster"}}

	config, err := clientcmd.BuildConfigFromFlags("", opts.Kubeconfig)
	if err != nil {
		return nil, err
	}

	collector.clientset, err = dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return collector, nil
}

func (c *ClusterCollector) Get() ([]interface{}, error) {
	gvrs := []schema.GroupVersionResource{
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"},
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"},
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"},
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
		schema.GroupVersionResource{Group: "policy", Version: "v1beta1", Resource: "podsecuritypolicies"},
		schema.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"},
	}

	var results []interface{}
	for _, g := range gvrs {
		ri := c.clientset.Resource(g)
		rs, err := ri.List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var resource interface{}

		for _, r := range rs.Items {
			if jsonManifest, ok := r.GetAnnotations()["kubectl.kubernetes.io/last-applied-configuration"]; ok {

				err := json.Unmarshal([]byte(jsonManifest), &resource)
				if err != nil {
					return nil, err
				}
				results = append(results, resource)
			}
		}
	}

	return results, nil
}
