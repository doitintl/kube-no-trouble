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

func newKubeCollector(kubeconfig string, kubecontext string, discoveryClient discovery.DiscoveryInterface, userAgent string) (*kubeCollector, error) {
	col := &kubeCollector{}

	if discoveryClient != nil {
		col.discoveryClient = discoveryClient
	} else {
		var err error
		if col.restConfig, err = newClientRestConfig(kubeconfig, kubecontext, rest.InClusterConfig, userAgent); err != nil {
			return nil, fmt.Errorf("failed to assemble client config: %w", err)
		}

		if col.discoveryClient, err = discovery.NewDiscoveryClientForConfig(col.restConfig); err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
	}

	return col, nil
}

func newClientRestConfig(kubeconfig string, kubecontext string, inClusterFn func() (*rest.Config, error), userAgent string) (*rest.Config, error) {
	if kubeconfig == "" {
		if restConfig, err := inClusterFn(); err == nil {
			restConfig.UserAgent = userAgent
			restConfig.WarningHandler = rest.NoWarnings{}
			return restConfig, nil
		}
	}

	pathOptions := clientcmd.NewDefaultPathOptions()
	if kubeconfig != "" {
		pathOptions.GlobalFile = kubeconfig
	}

	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	configOverrides := clientcmd.ConfigOverrides{}
	if kubecontext != "" {
		configOverrides.CurrentContext = kubecontext
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &configOverrides)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	restConfig.UserAgent = userAgent
	restConfig.WarningHandler = rest.NoWarnings{}
	return restConfig, nil
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
