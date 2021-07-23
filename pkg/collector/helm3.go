package collector

import (
	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/releaseutil"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/discovery"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type HelmV3Collector struct {
	*commonCollector
	*kubeCollector
	client       corev1.CoreV1Interface
	secretsStore *storage.Storage
	configStore  *storage.Storage
}

type HelmV3Opts struct {
	Kubeconfig      string
	KubeContext     string
	DiscoveryClient discovery.DiscoveryInterface
	CoreClient      corev1.CoreV1Interface
}

func NewHelmV3Collector(opts *HelmV3Opts) (*HelmV3Collector, error) {
	kubeCollector, err := newKubeCollector(opts.Kubeconfig, opts.KubeContext, opts.DiscoveryClient)
	if err != nil {
		return nil, err
	}

	collector := &HelmV3Collector{
		commonCollector: newCommonCollector("Helm v3"),
		kubeCollector:   kubeCollector,
	}

	if opts.CoreClient != nil {
		collector.client = opts.CoreClient
	} else if collector.client, err = corev1.NewForConfig(kubeCollector.GetRestConfig()); err != nil {
		return nil, err
	}

	secretsDriver := driver.NewSecrets(collector.client.Secrets(""))
	collector.secretsStore = storage.Init(secretsDriver)

	configDriver := driver.NewConfigMaps(collector.client.ConfigMaps(""))
	collector.configStore = storage.Init(configDriver)

	return collector, nil
}

func (c *HelmV3Collector) Get() ([]map[string]interface{}, error) {
	releases, err := c.secretsStore.ListDeployed()
	if err != nil {
		return nil, err
	}

	releasesConfig, err := c.configStore.ListDeployed()
	if err != nil {
		return nil, err
	}

	releases = append(releases, releasesConfig...)

	var results []map[string]interface{}

	for _, r := range releases {
		manifests := releaseutil.SplitManifests(r.Manifest)
		for _, m := range manifests {
			var manifest map[string]interface{}

			err := yaml.Unmarshal([]byte(m), &manifest)
			if err != nil {
				log.Warn().Msgf("failed to parse release %s/%s: %v", r.Namespace, r.Name, err)
				continue
			}

			// Default to the release namespace if the manifest doesn't have the namespace set
			if meta, ok := manifest["metadata"]; ok {
				switch v := meta.(type) {
				case map[string]interface{}:
					if _, ok := v["namespace"]; !ok {
						v["namespace"] = r.Namespace
					}
				}
			}

			results = append(results, manifest)
		}
	}

	return results, nil
}
