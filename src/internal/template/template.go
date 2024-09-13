package template

import (
	"os"
	"strings"
	"text/template"

	"github.com/defenseunicorns/pkg/helpers"
)

const PREFIX = "LULA_"

// ExecuteTemplate templates the template string with the data map
func ExecuteTemplate(data map[string]interface{}, templateString string) ([]byte, error) {
	tmpl, err := template.New("template").Parse(templateString)
	if err != nil {
		return []byte{}, err
	}
	tmpl.Option("missingkey=default")

	var buffer strings.Builder
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// Prepare the map of data for use in templating

func CollectTemplatingData(data map[string]interface{}) map[string]interface{} {

	// Get all environment variables with a specific prefix
	envMap := GetEnvVars(PREFIX)

	// Merge the data into a single map for use with templating
	mergedMap := helpers.MergeMapRecursive(envMap, data)

	return mergedMap

}

// get all environment variables with the established prefix
func GetEnvVars(prefix string) map[string]interface{} {
	envMap := make(map[string]interface{})

	// Get all environment variables
	envVars := os.Environ()

	// Iterate over environment variables
	for _, envVar := range envVars {
		// Split the environment variable into key and value
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		// Check if the key starts with the specified prefix
		if strings.HasPrefix(key, prefix) {
			// Remove the prefix from the key and convert to lowercase
			strippedKey := strings.TrimPrefix(key, prefix)
			envMap[strings.ToLower(strippedKey)] = value
		}
	}

	return envMap
}
