package collector

import (
	"github.com/rs/zerolog/log"
	"helm.sh/helm/pkg/storage"
	"helm.sh/helm/pkg/storage/driver"
	"k8s.io/client-go/discovery"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const HELM_V2_COLLECTOR_NAME = "Helm v2"

type HelmV2Collector struct {
	*commonCollector
	*kubeCollector
	client       corev1.CoreV1Interface
	secretsStore *storage.Storage
	configStore  *storage.Storage
}

type HelmV2Opts struct {
	Kubeconfig      string
	KubeContext     string
	DiscoveryClient discovery.DiscoveryInterface
	CoreClient      corev1.CoreV1Interface
}

func NewHelmV2Collector(opts *HelmV2Opts, userAgent string) (*HelmV2Collector, error) {

	kubeCollector, err := newKubeCollector(opts.Kubeconfig, opts.KubeContext, opts.DiscoveryClient, userAgent)
	if err != nil {
		return nil, err
	}

	collector := &HelmV2Collector{
		commonCollector: newCommonCollector(HELM_V2_COLLECTOR_NAME),
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

func (c *HelmV2Collector) Get() ([]map[string]interface{}, error) {
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
		if manifests, err := parseManifests(r.Manifest, r.Namespace); err != nil {
			log.Warn().Msgf("failed to parse release %s/%s: %v", r.Namespace, r.Name, err)
		} else {
			results = append(results, manifests...)
		}
	}

	return results, nil
}
