package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

func ParseTemplateOverrides(setFlags []string) (map[string]string, error) {
	overrides := make(map[string]string)
	for _, flag := range setFlags {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			return overrides, fmt.Errorf("invalid --set flag format, should be .root.key=value")
		}

		if !strings.HasPrefix(parts[0], "."+template.CONST+".") && !strings.HasPrefix(parts[0], "."+template.VAR+".") {
			return overrides, fmt.Errorf("invalid --set flag format, path should start with .const or .var")
		}

		path, value := parts[0], parts[1]
		overrides[path] = value
	}
	return overrides, nil
}

// writeResources writes the resources to a file or stdout
func WriteResources(data types.DomainResources, filepath string) error {
	jsonData := message.JSONValue(data)

	// If a filepath is provided, write the JSON data to the file.
	if filepath != "" {
		err := os.WriteFile(filepath, []byte(jsonData), 0600)
		if err != nil {
			return fmt.Errorf("error writing resource JSON to file: %v", err)
		}
	} else {
		message.Printf("%s", jsonData)
	}
	return nil
}
