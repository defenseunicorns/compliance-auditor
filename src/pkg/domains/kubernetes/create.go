package kube

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/message"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
)

// CreateAllResources() creates all resources and returns their status
func CreateAllResources(ctx context.Context, cluster *Cluster, resources []CreateResource) (map[string]interface{}, []string, error) {
	collections := make(map[string]interface{}, len(resources))
	namespaces := make([]string, 0)
	var errList []string

	if cluster == nil {
		return nil, nil, fmt.Errorf("cluster is nil")
	}

	// Create the resources, collect the outcome
	for _, resource := range resources {
		var collection []map[string]interface{}
		var err error
		// Create namespace if specified
		if resource.Namespace != "" {
			new, err := createNamespace(ctx, cluster.kclient, resource.Namespace)
			if err != nil {
				message.Debugf("error creating namespace %s: %v", resource.Namespace, err)
				errList = append(errList, err.Error())
			}
			// Only add to list of resources to clean up if not already in cluster
			if new {
				namespaces = append(namespaces, resource.Namespace)
			}
		}

		// TODO: Allow both Manifest and File to be specified?
		// Want to catch any errors and proceed in case resources have already been created
		if resource.Manifest != "" {
			collection, err = CreateFromManifest(ctx, cluster.kclient, []byte(resource.Manifest))
			if err != nil {
				message.Debugf("error creating resource from manifest: %v", err)
				errList = append(errList, err.Error())
			}
		} else if resource.File != "" {
			collection, err = CreateFromFile(ctx, cluster.kclient, resource.File)
			if err != nil {
				message.Debugf("error creating resource from file: %v", err)
				errList = append(errList, err.Error())
			}
		} else {
			errList = append(errList, "resource must have either manifest or file specified")
		}
		collections[resource.Name] = collection
	}

	// Check if there were any errors
	if len(errList) > 0 {
		return nil, nil, errors.New("errors creating resources encountered: " + strings.Join(errList, "; "))
	}

	return collections, namespaces, nil
}

// CreateResourceFromManifest() creates the resource from the manifest string
func CreateFromManifest(ctx context.Context, client klient.Client, resourceBytes []byte) ([]map[string]interface{}, error) {
	resources := make([]map[string]interface{}, 0)

	objArray, err := readResourcesFromYaml(resourceBytes)
	if err != nil {
		return nil, err
	}
	for _, obj := range objArray {
		resource, err := createResource(ctx, client, &obj)
		if err == nil {
			resources = append(resources, resource.Object)
		}
	}

	cleanResources(&resources)

	return resources, nil
}

// CreateResourceFromFile() creates the resource from a file
func CreateFromFile(ctx context.Context, client klient.Client, resourceFile string) ([]map[string]interface{}, error) {
	// Get manifest data from file and pass to CreateFromManifest
	resourceBytes, err := network.Fetch(resourceFile)
	if err != nil {
		return nil, err
	}
	return CreateFromManifest(ctx, client, resourceBytes)
}

// DestroyAllResources() removes all the created resources
func DestroyAllResources(ctx context.Context, client klient.Client, collections map[string]interface{}, namespaces []string) error {
	var errList []string // Collect errors to return at end so all resources are attempted to be destroyed
	for _, resources := range collections {
		if resources, ok := resources.([]map[string]interface{}); ok {
			// Destroy in reverse order
			for i := len(resources) - 1; i >= 0; i-- {
				obj := &unstructured.Unstructured{Object: resources[i]}
				err := destroyResource(ctx, client, obj)
				if err != nil {
					message.Debugf("error destroying resource %s: %v", obj.GetName(), err)
					errList = append(errList, err.Error())
				}
			}
		}
	}

	// Delete namespaces
	for _, namespace := range namespaces {
		ns := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": namespace,
				},
			},
		}

		if err := destroyResource(ctx, client, ns); err != nil {
			message.Debugf("error destroying namespace %s: %v", namespace, err)
			errList = append(errList, err.Error())
		}
	}

	// Check if there were any errors
	if len(errList) > 0 {
		return errors.New("errors encountered destroying resources: " + strings.Join(errList, "; "))
	}

	return nil
}

// createResource() creates a resource in a k8s cluster
func createResource(ctx context.Context, client klient.Client, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	// Create the object -> error returned when object is unable to be created
	if err := client.Resources().Create(ctx, obj); err != nil {
		return nil, err
	}

	// Wait for object to exist -> Times out at 30 seconds
	conditionFunc := func(obj k8s.Object) bool {
		if err := client.Resources().Get(ctx, obj.GetName(), obj.GetNamespace(), obj); err != nil {
			return false
		}
		return true
	}
	if err := wait.For(
		conditions.New(client.Resources()).ResourceMatch(obj, conditionFunc),
		wait.WithTimeout(time.Second*30),
	); err != nil {
		return nil, nil // Not returning error, just assuming that the object was blocked or not created
	}

	// Add pause for resources to do thier thang -> this should be subsumed by the addition of wait and resources
	time.Sleep(time.Second * 2) // Not sure if this is enough time, need to test with more complex resources

	// Get the object to return
	if err := client.Resources().Get(ctx, obj.GetName(), obj.GetNamespace(), obj); err != nil {
		return nil, err // Object was unable to be retrieved
	}

	return obj, nil
}

// destroyResource() removes a resource from a k8s cluster
func destroyResource(ctx context.Context, client klient.Client, obj *unstructured.Unstructured) error {
	propagationPolicy := metav1.DeletePropagationForeground
	if err := client.Resources().Delete(ctx, obj, resources.WithDeletePropagation(string(propagationPolicy))); err != nil {
		return err
	}

	// Wait for object to be removed from the cluster -> Times out at 5 minutes
	if err := wait.For(
		conditions.New(client.Resources()).ResourceDeleted(obj),
		wait.WithTimeout(time.Minute*5),
	); err != nil {
		return err // Object is unable to be deleted... retry logic? Or just return error?
	}

	return nil
}

// createNamespace() creates a namespace in a k8s cluster
func createNamespace(ctx context.Context, client klient.Client, namespace string) (new bool, err error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	if err := client.Resources().Get(ctx, namespace, "", ns); err == nil {
		return false, nil // Namespace already exists
	}

	if err := client.Resources().Create(ctx, ns); err != nil {
		return false, err // Namespace was unable to be created
	}

	return true, nil // Namespace created successfully
}

// readResourcesFromYaml reads a yaml file of k8s resources to an array of resources
func readResourcesFromYaml(resourceBytes []byte) (resources []unstructured.Unstructured, err error) {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(resourceBytes), 4096)

	for {
		resource := &unstructured.Unstructured{}
		if err := decoder.Decode(resource); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		resources = append(resources, *resource)
	}

	return resources, nil
}
