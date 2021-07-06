package collector

import (
	"fmt"

	goversion "github.com/hashicorp/go-version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

type kubeCollector struct {
	discoveryClient discovery.DiscoveryInterface
}

func newKubeCollector(kubeconfig string, discoveryClient discovery.DiscoveryInterface) (*kubeCollector, error) {
	col := &kubeCollector{}

	if discoveryClient == nil {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}

		col.discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
		if err != nil {
			return nil, err
		}
	} else {
		col.discoveryClient = discoveryClient
	}

	return col, nil
}

func (c *kubeCollector) GetServerVersion() (*goversion.Version, error) {
	version, err := c.discoveryClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version %w", err)
	}

	return goversion.NewVersion(version.String())
}
