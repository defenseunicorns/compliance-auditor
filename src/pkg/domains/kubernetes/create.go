package kube

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
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

// CreateE2E() creates the test resources, reads status, and destroys them
func CreateE2E(ctx context.Context, resources []CreateResource) (map[string]interface{}, error) {
	collections := make(map[string]interface{}, len(resources))
	namespaces := make([]string, 0)

	// Set up the clients
	config, err := connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to k8s cluster: %w", err)
	}
	client, err := klient.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create e2e client: %w", err)
	}

	// Create the resources, collect the outcome
	for _, resource := range resources {
		var collection []map[string]interface{}
		var err error
		// Create namespace if specified
		if resource.Namespace != "" {
			new, err := createNamespace(ctx, client, resource.Namespace)
			if err != nil {
				return nil, err
			}
			// Only add to list if not already in cluster
			if new {
				namespaces = append(namespaces, resource.Namespace)
			}
		}

		// TODO: Allow both Manifest and File to be specified?
		if resource.Manifest != "" {
			collection, err = CreateFromManifest(ctx, client, []byte(resource.Manifest))
			if err != nil {
				return nil, err
			}
		} else if resource.File != "" {
			collection, err = CreateFromFile(ctx, client, resource.File)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("resource must have either manifest or file specified")
		}
		collections[resource.Name] = collection
	}

	// Destroy the resources
	if err = DestroyAllResources(ctx, client, collections, namespaces); err != nil {
		// What do you do if resources can't be destroyed...
		return nil, err
	}
	return collections, nil
}

// CreateResourceFromManifest() creates the resource from the manifest string
func CreateFromManifest(ctx context.Context, client klient.Client, resourceBytes []byte) ([]map[string]interface{}, error) {
	resources := make([]map[string]interface{}, 0)

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(resourceBytes), 4096)
	for {
		rawObj := &unstructured.Unstructured{}
		if err := decoder.Decode(rawObj); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		resource, err := createResource(ctx, client, rawObj)
		if err == nil {
			resources = append(resources, resource.Object)
		}
	}
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
	for _, resources := range collections {
		if resources, ok := resources.([]map[string]interface{}); ok {
			// Destroy in reverse order
			for i := len(resources) - 1; i >= 0; i-- {
				obj := &unstructured.Unstructured{Object: resources[i]}
				err := destroyResource(ctx, client, obj)
				if err != nil {
					return err // Should I try again if this returns an error, or handle some other way?
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
			return err
		}
	}

	return nil
}

// createResource() creates a resource in a k8s cluster
func createResource(ctx context.Context, client klient.Client, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	// Modify the obj name to avoid collisions
	//obj.SetName(envconf.RandomName(obj.GetName(), 16)) // Maybe don't do this... What if you want to check a specific object?

	// Create the object
	if err := client.Resources().Create(ctx, obj); err != nil {
		return nil, err
	}

	// Wait for object to exist -> Times out at 10 seconds -> Presumably the object is blocked(?)
	conditionFunc := func(obj k8s.Object) bool {
		if err := client.Resources().Get(ctx, obj.GetName(), obj.GetNamespace(), obj); err != nil {
			return false
		}
		return true
	}
	if err := wait.For(
		conditions.New(client.Resources()).ResourceMatch(obj, conditionFunc),
		wait.WithTimeout(time.Second*10),
	); err != nil {
		return nil, nil // Not returning error, just assuming that the object was blocked or not created
	}

	// Add pause for resources to do thier thang
	time.Sleep(time.Second * 2) // Not sure if this is enough time

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
