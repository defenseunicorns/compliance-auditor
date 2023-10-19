package kube

import (
	"context"

	"github.com/defenseunicorns/lula/src/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

// QueryCluster() requires context and a Payload as input and returns []unstructured.Unstructured
// This function is used to query the cluster for all resources required for processing
func QueryCluster(ctx context.Context, payload types.Payload) ([]unstructured.Unstructured, error) {

	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)
	var resources []unstructured.Unstructured

	for _, rule := range payload.ResourceRules {

		if len(rule.Namespaces) == 0 {
			items, err := GetResourcesDynamically(dynamic, ctx,
				rule.Group, rule.Version, rule.Resource, "")
			if err != nil {
				return nil, err
			}
			resources = append(resources, items...)
		} else {
			for _, namespace := range rule.Namespaces {
				items, err := GetResourcesDynamically(dynamic, ctx,
					rule.Group, rule.Version, rule.Resource, namespace)
				if err != nil {
					return nil, err
				}
				resources = append(resources, items...)
			}

		}

	}

	return resources, nil
}

// GetResourcesDynamically() requires a dynamic interface and processes GVR to return []unstructured.Unstructured
// This function is used to query the cluster for specific subset of resources required for processing
func GetResourcesDynamically(dynamic dynamic.Interface, ctx context.Context,
	group string, version string, resource string, namespace string) (
	[]unstructured.Unstructured, error) {

	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	list, err := dynamic.Resource(resourceId).Namespace(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
