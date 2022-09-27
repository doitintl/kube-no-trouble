package collector

import (
	"fmt"

	"helm.sh/helm/v3/pkg/releaseutil"
	"sigs.k8s.io/yaml"
)

func parseManifests(manifest string, defaultNamespace string) ([]MetaOject, error) {
	var results []MetaOject

	manifests := releaseutil.SplitManifests(manifest)
	for i, m := range manifests {
		var manifest MetaOject

		err := yaml.Unmarshal([]byte(m), &manifest)
		if err != nil {
			err = fmt.Errorf("failed to parse manifest %s: %v", i, err)
			return nil, err
		}

		fixNamespace(&manifest, defaultNamespace)

		results = append(results, manifest)
	}

	return results, nil
}

func fixNamespace(manifest *MetaOject, defaultNamespace string) {
	// Default to the release namespace if the manifest doesn't have the namespace set
	if manifest.Namespace == "" {
		manifest.Namespace = defaultNamespace
	}
}
