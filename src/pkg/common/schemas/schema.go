package schemas

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/model"
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed *.json
var Schemas embed.FS

const (
	SCHEMA_SUFFIX = ".json"
)

func PrefixSchema(path string) string {
	if !strings.HasSuffix(path, SCHEMA_SUFFIX) {
		path = path + SCHEMA_SUFFIX
	}
	return path
}

// HasSchema checks if a schema exists in the schemas directory
func HasSchema(path string) bool {
	path = PrefixSchema(path)
	_, err := Schemas.Open(path)
	return err == nil
}

// ListSchemas returns a list of schema names
func ListSchemas() ([]string, error) {
	files, err := ToMap()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	return keys, nil
}

// ToMap returns a map of schema names to schemas
func ToMap() (fileMap map[string]fs.DirEntry, err error) {
	files, err := Schemas.ReadDir(".")
	if err != nil {
		return nil, err
	}
	fileMap = make(map[string]fs.DirEntry)
	for _, file := range files {
		name := file.Name()
		isDir := file.IsDir()
		if isDir || !strings.HasSuffix(name, SCHEMA_SUFFIX) {
			continue
		}
		fileMap[name] = file
	}
	return fileMap, nil
}

// GetSchema returns a schema from the schemas directory
func GetSchema(path string) ([]byte, error) {
	path = PrefixSchema(path)
	if !HasSchema(path) {
		return nil, fmt.Errorf("schema not found")
	}
	return Schemas.ReadFile(path)
}

func Validate(schema string, data model.InterfaceOrBytes) error {
	jsonMap, err := model.CoerceToJsonMap(data)
	if err != nil {
		return err
	}

	schemaBytes, err := GetSchema(schema)
	if err != nil {
		return err
	}

	sch, err := jsonschema.CompileString(schema, string(schemaBytes))
	if err != nil {
		return err
	}

	err = sch.Validate(jsonMap)
	if err != nil {
		// If the error is not a validation error, return the error
		validationErr, ok := err.(*jsonschema.ValidationError)
		if !ok {
			return err
		}

		// Extract the specific errors from the schema error
		// Return the errors as a string
		basicOutput := validationErr.BasicOutput()
		basicErrors := validation.ExtractErrors(jsonMap, basicOutput)
		formattedErrors, _ := json.MarshalIndent(basicErrors, "", "  ")
		return errors.New(string(formattedErrors))
	}
	return nil
}
