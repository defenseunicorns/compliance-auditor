package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

// QueryCluster() requires context and a Payload as input and returns []unstructured.Unstructured
// This function is used to query the cluster for all resources required for processing
func QueryCluster(ctx context.Context, resources []types.Resource) (map[string]interface{}, error) {

	// We may need a new type here to hold groups of resources

	collections := make(map[string]interface{}, 0)

	for _, resource := range resources {
		collection, err := GetResourcesDynamically(ctx, resource.ResourceRule)
		// log error but continue with other resources
		if err != nil {
			return nil, err
		}

		if len(collection) > 0 {
			// Append to collections if not empty collection
			// convert to object if named resource
			if resource.ResourceRule.Name != "" {
				collections[resource.Name] = collection[0]
			} else {
				collections[resource.Name] = collection
			}
		}
	}
	return collections, nil
}

// GetResourcesDynamically() requires a dynamic interface and processes GVR to return []map[string]interface{}
// This function is used to query the cluster for specific subset of resources required for processing
func GetResourcesDynamically(ctx context.Context,
	resource types.ResourceRule) (
	[]map[string]interface{}, error) {

	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error with connection to the Cluster")
	}
	dynamic := dynamic.NewForConfigOrDie(config)

	resourceId := schema.GroupVersionResource{
		Group:    resource.Group,
		Version:  resource.Version,
		Resource: resource.Resource,
	}
	collection := make([]map[string]interface{}, 0)

	namespaces := []string{""}
	if len(resource.Namespaces) != 0 {
		namespaces = resource.Namespaces
	}
	for _, namespace := range namespaces {
		list, err := dynamic.Resource(resourceId).Namespace(namespace).
			List(ctx, metav1.ListOptions{})

		if err != nil {
			return nil, err
		}

		// Reduce if named resource
		if resource.Name != "" {
			// requires single specified namespace
			if len(resource.Namespaces) == 1 {
				item, err := reduceByName(resource.Name, list.Items)
				if err != nil {
					return nil, err
				}
				collection = append(collection, item)
				return collection, nil
			}

		} else {
			for _, item := range list.Items {
				collection = append(collection, item.Object)
			}
		}
	}

	return collection, nil
}

func getGroupVersionResource(kind string) (gvr *schema.GroupVersionResource, err error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error with connection to the Cluster")
	}
	name := strings.Split(kind, "/")[0]

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	_, resourceList, _, err := discoveryClient.GroupsAndMaybeResources()
	if err != nil {

		return nil, err
	}

	for gv, list := range resourceList {
		for _, item := range list.APIResources {
			if item.SingularName == name {
				return &schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: item.Name,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("kind %s not found", kind)
}

// reduceByName() takes a name and loops over all items to return the first match
func reduceByName(name string, items []unstructured.Unstructured) (map[string]interface{}, error) {

	for _, item := range items {
		if item.GetName() == name {
			return item.Object, nil
		}
	}

	return nil, fmt.Errorf("no resource found with name %s", name)
}
