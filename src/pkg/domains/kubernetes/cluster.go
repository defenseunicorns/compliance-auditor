package kube

import (
	"errors"
	"fmt"
	"sync"

	pkgkubernetes "github.com/defenseunicorns/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/cli-utils/pkg/kstatus/watcher"
	"sigs.k8s.io/e2e-framework/klient"
)

var (
	clusterConnectOnce  sync.Once
	globalCluster       *Cluster
	globalConnectionErr error
)

type Cluster struct {
	clientset     kubernetes.Interface
	kclient       klient.Client
	watcher       watcher.StatusWatcher
	dynamicClient *dynamic.DynamicClient
}

func GetCluster() (*Cluster, error) {
	clusterConnectOnce.Do(func() {
		globalCluster, globalConnectionErr = New()
	})

	return globalCluster, globalConnectionErr
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

	// Ensure no errors were returned to validate cluster connection.
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

	return nil, fmt.Errorf("resource %s not found in group, %s, version, %s", resource, group, version)
}
