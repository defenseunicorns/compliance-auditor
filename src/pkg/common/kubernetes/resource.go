package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/pkg/message"
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
			// Adding the collection to the map when empty will result in a false positive for the validation in OPA?
			// TODO: add warning log here
			collections[resource.Name] = collection
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
		return nil, fmt.Errorf("Error with connection to the Cluster")
	}
	dynamic := dynamic.NewForConfigOrDie(config)

	resourceId := schema.GroupVersionResource{
		Group:    resource.Group,
		Version:  resource.Version,
		Resource: resource.Resource,
	}
	collection := make([]map[string]interface{}, 0)
	// Query all namespaces in one execution
	if len(resource.Namespaces) == 0 {
		list, err := dynamic.Resource(resourceId).Namespace("").
			List(ctx, metav1.ListOptions{})

		if err != nil {
			return nil, err
		}

		// Reduce if named resource
		if resource.Name != "" {
			items, err := reduceByName(resource.Name, list.Items)
			if err != nil {
				return nil, err
			}
			for _, item := range items {
				collection = append(collection, item)
			}
		} else {
			for _, item := range list.Items {
				collection = append(collection, item.Object)
			}
		}
		// Query multiple namespaces
	} else {
		for _, namespace := range resource.Namespaces {
			list, err := dynamic.Resource(resourceId).Namespace(namespace).
				List(ctx, metav1.ListOptions{})

			if err != nil {
				return nil, err
			}

			// Reduce if named resource
			if resource.Name != "" {
				items, err := reduceByName(resource.Name, list.Items)
				if err != nil {
					return nil, err
				}
				for _, item := range items {
					collection = append(collection, item)
				}
			} else {
				for _, item := range list.Items {
					collection = append(collection, item.Object)
				}
			}
		}
	}

	return collection, nil
}

func getGroupVersionResource(kind string) (gvr *schema.GroupVersionResource, err error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("Error with connection to the Cluster")
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

// reduceByName takes a name and loops over objects to return all matches
// TODO: investigate supporting wildcard matching?
func reduceByName(name string, items []unstructured.Unstructured) ([]map[string]interface{}, error) {
	reducedItems := make([]map[string]interface{}, 0)
	for _, item := range items {
		if item.GetName() == name {
			reducedItems = append(reducedItems, item.Object)
		}
	}
	// TODO: Determine if this is an error or simply a log message
	if len(reducedItems) == 0 {
		message.Debugf("No resource found with Name %s \n", name)
	}
	return reducedItems, nil
}
