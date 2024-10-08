package kube

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/e2e-framework/klient"
)

var cluster *Cluster

type Cluster struct {
	config  *rest.Config
	kclient klient.Client
}

func InitCluster() error {
	if cluster != nil {
		return nil
	}

	c, err := New()
	if err == nil {
		cluster = c
	}
	return err
}

func New() (*Cluster, error) {
	config, err := connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to k8s cluster: %w", err)
	}
	kclient, err := klient.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create e2e client: %w", err)
	}

	return &Cluster{
		kclient: kclient,
		config:  config,
	}, nil
}

// Use the K8s "client-go" library to get the currently active kube context, in the same way that
// "kubectl" gets it if no extra config flags like "--kubeconfig" are passed.
func connect() (config *rest.Config, err error) {
	// Build the config from the currently active kube context in the default way that the k8s client-go gets it, which
	// is to look at the KUBECONFIG env var
	config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{}).ClientConfig()

	if err != nil {
		return nil, err
	}

	return config, nil
}
