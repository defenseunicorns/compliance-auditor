package types

import (
	"fmt"
	"os"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

type DomainResources map[string]interface{}

// WriteResources writes the domain resources to a file or stdout
func WriteResources(data DomainResources, filepath string) error {
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
