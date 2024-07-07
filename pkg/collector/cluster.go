package collector

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

const CLUSTER_COLLECTOR_NAME = "Cluster"

type ClusterCollector struct {
	*commonCollector
	*kubeCollector
	clientSet             dynamic.Interface
	additionalResources   []schema.GroupVersionResource
	additionalAnnotations []string
}

type ClusterOpts struct {
	Kubeconfig      string
	KubeContext     string
	ClientSet       dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
}

func NewClusterCollector(opts *ClusterOpts, additionalKinds []string, additionalAnnotations []string, userAgent string) (
	*ClusterCollector, error) {
	kubeCollector, err := newKubeCollector(opts.Kubeconfig, opts.KubeContext, opts.DiscoveryClient, userAgent)
	if err != nil {
		return nil, err
	}

	collector := &ClusterCollector{
		kubeCollector:         kubeCollector,
		commonCollector:       newCommonCollector(CLUSTER_COLLECTOR_NAME),
		additionalAnnotations: additionalAnnotations,
	}

	if opts.ClientSet == nil {
		collector.clientSet, err = dynamic.NewForConfig(kubeCollector.GetRestConfig())
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
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "replicasets"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
		{Group: "policy", Version: "v1beta1", Resource: "podsecuritypolicies"},
		{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
		{Group: "networking.k8s.io", Version: "v1", Resource: "ingressclasses"},
		{Group: "storage.k8s.io", Version: "v1", Resource: "csidrivers"},
		{Group: "storage.k8s.io", Version: "v1", Resource: "csinodes"},
		{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
		{Group: "storage.k8s.io", Version: "v1", Resource: "volumeattachments"},
		{Group: "storage.k8s.io", Version: "v1", Resource: "csistoragecapacities"},
		{Group: "scheduling.k8s.io", Version: "v1", Resource: "priorityclasses"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
		{Group: "coordination.k8s.io", Version: "v1", Resource: "leases"},
		{Group: "authorization.k8s.io", Version: "v1", Resource: "subjectaccessreviews"},
		{Group: "authorization.k8s.io", Version: "v1", Resource: "selfsubjectaccessreviews"},
		{Group: "authorization.k8s.io", Version: "v1", Resource: "localsubjectaccessreviews"},
		{Group: "authentication.k8s.io", Version: "v1", Resource: "tokenreviews"},
		{Group: "certificates.k8s.io", Version: "v1", Resource: "certificatesigningrequests"},
		{Group: "apiregistration.k8s.io", Version: "v1", Resource: "apiservices"},
		{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"},
		{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "mutatingwebhookconfigurations"},
		{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingwebhookconfigurations"},
		{Group: "node.k8s.io", Version: "v1", Resource: "runtimeclasses"},
		{Group: "policy", Version: "v1", Resource: "poddisruptionbudgets"},
		{Group: "policy", Version: "v1beta1", Resource: "podsecuritypolicies"},
		{Group: "discovery.k8s.io", Version: "v1", Resource: "endpointslices"},
		{Group: "batch", Version: "v1", Resource: "cronjobs"},
		{Group: "autoscaling", Version: "v2", Resource: "horizontalpodautoscalers"},
		{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"},
		{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshotclasses"},
		{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshotcontents"},
		{Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta2", Resource: "flowschemas"},
		{Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta2", Resource: "prioritylevelconfigurations"},
		{Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta3", Resource: "flowschemas"},
		{Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta3", Resource: "prioritylevelconfigurations"},
	}
	gvrs = append(gvrs, c.additionalResources...)

	var results []map[string]interface{}
	for _, g := range gvrs {
		ri := c.clientSet.Resource(g)
		log.Debug().Msgf("Retrieving: %s.%s.%s", g.Resource, g.Version, g.Group)
		rs, err := ri.List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Debug().Msgf("Failed to retrieve: %s: %s", g, err)
			continue
		}

		for _, r := range rs.Items {
			if jsonManifest, ok := c.getLastAppliedConfig(r.GetAnnotations()); ok {
				var manifest map[string]interface{}

				err := json.Unmarshal([]byte(jsonManifest), &manifest)
				if err != nil {
					log.Warn().Msgf("failed to parse 'last-applied-configuration' annotation of resource %s/%s: %v", r.GetNamespace(), r.GetName(), err)
					continue
				}
				results = append(results, manifest)
			}
		}
	}

	return results, nil
}

func (c *ClusterCollector) getLastAppliedConfig(resourceAnnotations map[string]string) (string, bool) {
	annotations := append([]string{"kubectl.kubernetes.io/last-applied-configuration"}, c.additionalAnnotations...)
	for _, annotation := range annotations {
		if jsonManifest, ok := resourceAnnotations[annotation]; ok {
			return jsonManifest, ok
		}
	}

	return "", false
}
