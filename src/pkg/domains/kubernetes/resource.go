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
)

// QueryCluster() requires context and a Payload as input and returns []unstructured.Unstructured
// This function is used to query the cluster for all resources required for processing
func QueryCluster(ctx context.Context, cluster *Cluster, resources []Resource) (map[string]interface{}, error) {
	if cluster == nil {
		return nil, fmt.Errorf("cluster is nil")
	}

	// We may need a new type here to hold groups of resources

	collections := make(map[string]interface{}, 0)

	for _, resource := range resources {
		collection, err := GetResourcesDynamically(ctx, cluster, resource.ResourceRule)
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
func GetResourcesDynamically(ctx context.Context, cluster *Cluster, resource *ResourceRule) ([]map[string]interface{}, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource rule is nil")
	}

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
		itemObj, err := cluster.dynamicClient.Resource(resourceId).Namespace(namespaces[0]).Get(ctx, resource.Name, metav1.GetOptions{})
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
			list, err := cluster.dynamicClient.Resource(resourceId).Namespace(namespace).
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
