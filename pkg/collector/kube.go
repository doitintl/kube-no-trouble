package collector

import (
	"fmt"
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

func (c *kubeCollector) GetServerVersion() (string, error) {
	version, err := c.discoveryClient.ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get server version %w", err)
	}
	return version.Major + "." + version.Minor, nil
}
