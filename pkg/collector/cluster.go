package collector

import (
	"encoding/json"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

type ClusterCollector struct {
	*commonCollector
	*kubeCollector
	clientSet           dynamic.Interface
	additionalResources []schema.GroupVersionResource
}

type ClusterOpts struct {
	Kubeconfig      string
	ClientSet       dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
}

func NewClusterCollector(opts *ClusterOpts, additionalKinds []string) (*ClusterCollector, error) {
	kubeCollector, err := newKubeCollector(opts.Kubeconfig, opts.DiscoveryClient)
	if err != nil {
		return nil, err
	}

	collector := &ClusterCollector{
		kubeCollector:   kubeCollector,
		commonCollector: newCommonCollector("Cluster"),
	}

	if opts.ClientSet == nil {
		config, err := clientcmd.BuildConfigFromFlags("", opts.Kubeconfig)
		if err != nil {
			return nil, err
		}

		collector.clientSet, err = dynamic.NewForConfig(config)
		if err != nil {
			return nil, err
		}

	} else {
		collector.clientSet = opts.ClientSet
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(collector.discoveryClient))
	for _, ar := range additionalKinds {
		gvk, _ := schema.ParseKindArg(ar)

		gvrMap, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Warn().Msgf("Failed to map %s Kind to resource: %s", gvk.Kind, err)
			continue
		}

		collector.additionalResources = append(collector.additionalResources, gvrMap.Resource)
	}

	return collector, nil
}

func (c *ClusterCollector) Get() ([]map[string]interface{}, error) {
	gvrs := []schema.GroupVersionResource{
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"},
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"},
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"},
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
		schema.GroupVersionResource{Group: "policy", Version: "v1beta1", Resource: "podsecuritypolicies"},
		schema.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingressclasses"},
		schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "csidrivers"},
		schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "csinodes"},
		schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
		schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "volumeattachments"},
		schema.GroupVersionResource{Group: "scheduling.k8s.io", Version: "v1", Resource: "priorityclasses"},
		schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
		schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
		schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
		schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
		schema.GroupVersionResource{Group: "coordination.k8s.io", Version: "v1", Resource: "leases"},
		schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "subjectaccessreviews"},
		schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "selfsubjectaccessreviews"},
		schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "localsubjectaccessreviews"},
		schema.GroupVersionResource{Group: "authentication.k8s.io", Version: "v1", Resource: "tokenreviews"},
		schema.GroupVersionResource{Group: "certificates.k8s.io", Version: "v1", Resource: "certificatesigningrequests"},
		schema.GroupVersionResource{Group: "apiregistration.k8s.io", Version: "v1", Resource: "apiservices"},
		schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"},
		schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "mutatingwebhookconfigurations"},
		schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingwebhookconfigurations"},
		schema.GroupVersionResource{Group: "node.k8s.io", Version: "v1", Resource: "runtimeclasses"},
		schema.GroupVersionResource{Group: "policy", Version: "v1", Resource: "poddisruptionbudgets"},
	}
	gvrs = append(gvrs, c.additionalResources...)

	var results []map[string]interface{}
	for _, g := range gvrs {
		ri := c.clientSet.Resource(g)
		log.Debug().Msgf("Retrieving: %s.%s.%s", g.Resource, g.Version, g.Group)
		rs, err := ri.List(metav1.ListOptions{})
		if err != nil {
			log.Warn().Msgf("Failed to retrieve: %s: %s", g, err)
			continue
		}

		for _, r := range rs.Items {
			if jsonManifest, ok := r.GetAnnotations()["kubectl.kubernetes.io/last-applied-configuration"]; ok {
				var manifest map[string]interface{}

				err := json.Unmarshal([]byte(jsonManifest), &manifest)
				if err != nil {
					err := fmt.Errorf("failed to parse 'last-applied-configuration' annotation of resource %s/%s: %v", r.GetNamespace(), r.GetName(), err)
					return nil, err
				}
				results = append(results, manifest)
			}
		}
	}

	return results, nil
}
