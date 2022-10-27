package collector

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/releaseutil"
)

func parseManifests(manifest string, defaultNamespace string, discoveryClient discovery.DiscoveryInterface) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	manifests := releaseutil.SplitManifests(manifest)
	for i, m := range manifests {
		var manifest map[string]interface{}

		err := yaml.Unmarshal([]byte(m), &manifest)
		if err != nil {
			err = fmt.Errorf("failed to parse manifest %s: %v", i, err)
			return nil, err
		}
		// the helm's SplitManifests can keep empty documents in some cases
		if len(manifest) != 0 {
			umanifest := &unstructured.Unstructured{Object: manifest}
			log.Debug().Msgf("retrieved: %s/%s (%s)", umanifest.GetNamespace(), umanifest.GetName(), umanifest.GroupVersionKind())

			fixNamespace(umanifest, defaultNamespace, discoveryClient)

			results = append(results, manifest)
		}
	}

	return results, nil
}

func fixNamespace(resource *unstructured.Unstructured, defaultNamespace string, discoveryClient discovery.DiscoveryInterface) {
	// Default to the release namespace if the manifest doesn't have the namespace set
	if resource.GetNamespace() == "" && isResourceNamespaced(discoveryClient, resource.GroupVersionKind()) {
		resource.SetNamespace(defaultNamespace)
	}
}

func isResourceNamespaced(discoveryClient discovery.DiscoveryInterface, gvk schema.GroupVersionKind) bool {
	rs, err := discoveryClient.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	// It seems discovery client fails with error if resource is not found
	// this can happen, but is not fatal, we should notify user but continue
	if err != nil {
		log.Warn().Msgf("failed to discover supported resources for %s: %v", gvk.GroupVersion(), err)
		return false
	}

	for _, r := range rs.APIResources {
		if r.Kind == gvk.Kind {
			return r.Namespaced
		}
	}

	return false
}
