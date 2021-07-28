package collector

import (
	"fmt"

	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/releaseutil"
)

func parseManifests(manifest string, defaultNamespace string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	manifests := releaseutil.SplitManifests(manifest)
	for i, m := range manifests {
		var manifest map[string]interface{}

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

func fixNamespace(manifest *map[string]interface{}, defaultNamespace string) {
	// Default to the release namespace if the manifest doesn't have the namespace set
	if meta, ok := (*manifest)["metadata"]; ok {
		switch v := meta.(type) {
		case map[string]interface{}:
			if val, ok := v["namespace"]; !ok || val == nil {
				v["namespace"] = defaultNamespace
			}
		}
	}
}
