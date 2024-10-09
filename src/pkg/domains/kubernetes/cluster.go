package kube

import (
	"errors"
	"fmt"

	pkgkubernetes "github.com/defenseunicorns/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/cli-utils/pkg/kstatus/watcher"
	"sigs.k8s.io/e2e-framework/klient"
)

var globalCluster *Cluster

type Cluster struct {
	clientset     kubernetes.Interface
	kclient       klient.Client
	watcher       watcher.StatusWatcher
	dynamicClient *dynamic.DynamicClient
}

func InitCluster() error {
	if globalCluster != nil {
		return nil
	}

	c, err := New()
	if err == nil {
		globalCluster = c
	}
	return err
}

func New() (*Cluster, error) {
	clusterErr := errors.New("unable to connect to the cluster")
	clientset, config, err := pkgkubernetes.ClientAndConfig()
	if err != nil {
		return nil, errors.Join(clusterErr, err)
	}

	watcher, err := pkgkubernetes.WatcherForConfig(config)
	if err != nil {
		return nil, errors.Join(clusterErr, err)
	}

	kclient, err := klient.New(config)
	if err != nil {
		return nil, errors.Join(clusterErr, err)
	}

	dynamicClient := dynamic.NewForConfigOrDie(config)

	// Dogsled the version output. We just want to ensure no errors were returned to validate cluster connection.
	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, errors.Join(clusterErr, err)
	}

	return &Cluster{
		clientset:     clientset,
		kclient:       kclient,
		watcher:       watcher,
		dynamicClient: dynamicClient,
	}, nil
}

func (c *Cluster) validateAndGetGVR(group, version, resource string) (*metav1.APIResource, error) {
	// Create a discovery client
	discoveryClient := c.clientset.Discovery()

	// Get a list of all API resources for the given group version
	gv := schema.GroupVersion{
		Group:   group,
		Version: version,
	}
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion(gv.String())
	if err != nil {
		return nil, err
	}

	// Search for the specified resource in the list
	for _, apiResource := range resourceList.APIResources {
		if apiResource.Name == resource {
			return &apiResource, nil
		}
	}

	return nil, fmt.Errorf("resource %s not found in group %s version %s", resource, group, version)
}

// Use the K8s "client-go" library to get the currently active kube context, in the same way that
// "kubectl" gets it if no extra config flags like "--kubeconfig" are passed.
// func connect() (config *rest.Config, err error) {
// 	// Build the config from the currently active kube context in the default way that the k8s client-go gets it, which
// 	// is to look at the KUBECONFIG env var
// 	config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
// 		clientcmd.NewDefaultClientConfigLoadingRules(),
// 		&clientcmd.ConfigOverrides{}).ClientConfig()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return config, nil
// }
