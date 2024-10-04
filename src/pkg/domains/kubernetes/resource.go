package kube

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// QueryCluster() requires context and a Payload as input and returns []unstructured.Unstructured
// This function is used to query the cluster for all resources required for processing
func QueryCluster(ctx context.Context, resources []Resource) (map[string]interface{}, error) {

	// We may need a new type here to hold groups of resources

	collections := make(map[string]interface{}, 0)

	for _, resource := range resources {
		collection, err := GetResourcesDynamically(ctx, resource.ResourceRule)
		// log error but continue with other resources
		if err != nil {
			return nil, err
		}

		if resource.ResourceRule.Name != "" {
			if len(collection) > 0 {
				collections[resource.Name] = collection[0]
			} else {
				// This request returned no resources
				collections[resource.Name] = map[string]interface{}{}
			}

		} else {
			if len(collection) > 0 {
				collections[resource.Name] = collection
			} else {
				// This request returned no resources
				collections[resource.Name] = []map[string]interface{}{}
			}
		}
	}
	return collections, nil
}

// GetResourcesDynamically() requires a dynamic interface and processes GVR to return []map[string]interface{}
// This function is used to query the cluster for specific subset of resources required for processing
func GetResourcesDynamically(ctx context.Context,
	resource *ResourceRule) (
	[]map[string]interface{}, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource rule is nil")
	}
	config, err := connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to k8s cluster: %w", err)
	}
	dynamic := dynamic.NewForConfigOrDie(config)

	resourceId := schema.GroupVersionResource{
		Group:    resource.Group,
		Version:  resource.Version,
		Resource: resource.Resource,
	}
	collection := make([]map[string]interface{}, 0)

	// Depending on resource-rule, either a single item or list of items will be appended to collection
	namespaces := resource.Namespaces
	if len(namespaces) == 0 {
		namespaces = []string{""}
	}

	if resource.Name != "" && len(namespaces) > 1 {
		return nil, fmt.Errorf("named resource requested cannot be returned from multiple namespaces")
	} else if resource.Name != "" {
		// Extracting named resources can only occur here
		var itemObj *unstructured.Unstructured
		itemObj, err := dynamic.Resource(resourceId).Namespace(namespaces[0]).Get(ctx, resource.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		item := itemObj.Object

		// If field is specified, get the field data; can only occur when a single named resource is specified
		if resource.Field != nil && resource.Field.Jsonpath != "" {
			item, err = getFieldValue(item, resource.Field)
			if err != nil {
				return nil, err
			}
		}

		collection = append(collection, item)
	} else {
		for _, namespace := range namespaces {
			list, err := dynamic.Resource(resourceId).Namespace(namespace).
				List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, err
			}

			for _, item := range list.Items {
				collection = append(collection, item.Object)
			}
		}
	}

	cleanResources(&collection)

	return collection, nil
}

func getGroupVersionResource(kind string) (gvr *schema.GroupVersionResource, err error) {
	config, err := connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to k8s cluster: %w", err)
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

// getFieldValue() looks up the field from a resource and returns a map[string]interface{} representation of the data
func getFieldValue(item map[string]interface{}, field *Field) (map[string]interface{}, error) {
	if field == nil {
		return nil, fmt.Errorf("field is nil")
	}
	// Identify the field in item
	pathParts := strings.Split(field.Jsonpath, ".")[1:]
	current := item
	var fieldValue string
	for i, part := range pathParts {
		// Check if the first part is a valid key
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			if value, ok := current[strings.Join(pathParts[i:], ".")].(string); ok {
				fieldValue = value
				break
			} else {
				return nil, fmt.Errorf("path not found: %s", strings.Join(pathParts[:i+1], "."))
			}

		}
	}

	// If base64 encoded, decode the data first
	if field.Base64 {
		decoded, err := base64.StdEncoding.DecodeString(fieldValue)
		if err != nil {
			return nil, err
		}
		fieldValue = string(decoded)
	}

	// If field type is unset, set to default
	if field.Type == "" {
		field.Type = DefaultFieldType
	}

	var data interface{}
	// Get the field data if json
	if field.Type == FieldTypeJSON {
		// Convert fieldValue to json
		err := json.Unmarshal([]byte(fieldValue), &data)
		if err != nil {
			return nil, err
		}
		result, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected JSON to decode for field %s", field.Jsonpath)
		}
		return result, nil
	} else {
		// Convert fieldValue to yaml
		err := yaml.Unmarshal([]byte(fieldValue), &data)
		if err != nil {
			return nil, err
		}
		result, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected YAML to decode field %s", field.Jsonpath)
		}
		return result, nil
	}
}

// cleanResources() clears out unnecceary fields from the resources that contribute to noise
func cleanResources(resources *[]map[string]interface{}) {
	// Removes metadata.managedFields from each item in the collection
	// Field is long and seemingly useless for our purposes, removing to reduce noise
	for _, c := range *resources {
		if metadata, ok := c["metadata"].(map[string]interface{}); ok {
			delete(metadata, "managedFields")
		}
	}
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
