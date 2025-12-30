package kube

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildRESTConfig(cfg *Config) (*rest.Config, error) {
	if cfg.KubeConfigPath != "" {
		loading := &clientcmd.ClientConfigLoadingRules{
			ExplicitPath: cfg.KubeConfigPath,
		}
		overrides := &clientcmd.ConfigOverrides{}
		if cfg.KubeContext != "" {
			overrides.CurrentContext = cfg.KubeContext
		}

		cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loading, overrides)
		restCfg, err := cc.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("kube -> load kubeconfig %q -> %w", cfg.KubeConfigPath, err)
		}
		return restCfg, nil
	}

	restCfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("kube -> in-cluster config -> %w", err)
	}
	return restCfg, nil
}
