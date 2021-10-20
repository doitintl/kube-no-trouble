package collector

import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type kubeCollector struct {
	discoveryClient discovery.DiscoveryInterface
	restConfig      *rest.Config
}

func newKubeCollector(kubeconfig string, kubecontext string, discoveryClient discovery.DiscoveryInterface) (*kubeCollector, error) {
	col := &kubeCollector{}

	if discoveryClient != nil {
		col.discoveryClient = discoveryClient
	} else {
		pathOptions := clientcmd.NewDefaultPathOptions()
		if kubeconfig != "" {
			pathOptions.GlobalFile = kubeconfig
		}

		config, err := pathOptions.GetStartingConfig()

		configOverrides := clientcmd.ConfigOverrides{}
		if kubecontext != "" {
			configOverrides.CurrentContext = kubecontext
		}

		clientConfig := clientcmd.NewDefaultClientConfig(*config, &configOverrides)
		col.restConfig, err = clientConfig.ClientConfig()
		if err != nil {
			return nil, err
		}

		if col.discoveryClient, err = discovery.NewDiscoveryClientForConfig(col.restConfig); err != nil {
			return nil, err
		}
	}

	return col, nil
}

func (c *kubeCollector) GetServerVersion() (*judge.Version, error) {
	version, err := c.discoveryClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version %w", err)
	}

	return judge.NewVersion(version.String())
}

func (c *kubeCollector) GetRestConfig() *rest.Config {
	return c.restConfig
}
