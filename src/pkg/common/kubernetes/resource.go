package kube

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/types"
	"gopkg.in/yaml.v3"
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
				// If field is specified, get the field data
				if resource.Field.Jsonpath != "" {
					item, err = getFieldValue(item, resource.Field)
					if err != nil {
						return nil, err
					}
				}

				collection = append(collection, item)
			}

		} else {
			for _, item := range list.Items {
				collection = append(collection, item.Object)
			}
		}
	}

	// Removes metadata.managedFields from each item in the collection
	// Field is long and seemingly useless for our purposes, removing to reduce noise
	for _, c := range collection {
		if metadata, ok := c["metadata"].(map[string]interface{}); ok {
			delete(metadata, "managedFields")
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

// getFieldValue() looks up the field from a resource and returns a map[string]interface{} representation of the data
func getFieldValue(item map[string]interface{}, field types.Field) (map[string]interface{}, error) {

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
		field.Type = types.DefaultFieldType
	}

	var data interface{}
	// Get the field data if json
	if field.Type == types.FieldTypeJSON {
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
