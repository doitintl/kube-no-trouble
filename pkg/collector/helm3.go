package collector

import (
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/discovery"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const HELM_V3_COLLECTOR_NAME = "Helm v3"

type HelmV3Collector struct {
	*commonCollector
	*kubeCollector
	client        corev1.CoreV1Interface
	secretsStores []*storage.Storage
	configStores  []*storage.Storage
}

type HelmV3Opts struct {
	Kubeconfig      string
	KubeContext     string
	DiscoveryClient discovery.DiscoveryInterface
	CoreClient      corev1.CoreV1Interface
}

func NewHelmV3Collector(opts *HelmV3Opts, namespaces []string, userAgent string) (*HelmV3Collector, error) {
	kubeCollector, err := newKubeCollector(opts.Kubeconfig, opts.KubeContext, opts.DiscoveryClient, userAgent)
	if err != nil {
		return nil, err
	}

	collector := &HelmV3Collector{
		commonCollector: newCommonCollector(HELM_V3_COLLECTOR_NAME),
		kubeCollector:   kubeCollector,
	}

	if opts.CoreClient != nil {
		collector.client = opts.CoreClient
	} else if collector.client, err = corev1.NewForConfig(kubeCollector.GetRestConfig()); err != nil {
		return nil, err
	}

	for _, namespace := range namespaces {
		secretsDriver := driver.NewSecrets(collector.client.Secrets(namespace))
		collector.secretsStores = append(collector.secretsStores, storage.Init(secretsDriver))

		configDriver := driver.NewConfigMaps(collector.client.ConfigMaps(namespace))
		collector.configStores = append(collector.configStores, storage.Init(configDriver))
	}

	return collector, nil
}

func (c *HelmV3Collector) Get() ([]map[string]interface{}, error) {
	var releases []*release.Release

	for _, secretsStore := range c.secretsStores {
		releasesSecret, err := secretsStore.ListDeployed()
		if err != nil {
			return nil, err
		}
		releases = append(releases, releasesSecret...)
	}

	for _, configStore := range c.configStores {
		releasesConfig, err := configStore.ListDeployed()
		if err != nil {
			return nil, err
		}
		releases = append(releases, releasesConfig...)
	}

	var results []map[string]interface{}

	for _, r := range releases {
		if manifests, err := parseManifests(r.Manifest, r.Namespace, c.discoveryClient); err != nil {
			log.Warn().Msgf("failed to parse release %s/%s: %v", r.Namespace, r.Name, err)
		} else {
			results = append(results, manifests...)
		}
	}

	return results, nil
}
