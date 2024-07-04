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
		basicErrors := ExtractErrors(jsonMap, basicOutput)
		formattedErrors, _ := json.MarshalIndent(basicErrors, "", "  ")
		return errors.New(string(formattedErrors))
	}
	return nil
}

// Creates a []ValidatorError from a jsonschema.Basic
// The jsonschema.Basic contains the errors from the validation
func ExtractErrors(originalObject map[string]interface{}, validationError jsonschema.Basic) (validationErrors []validation.ValidatorError) {
	validationErrors = []validation.ValidatorError{}
	for _, basicError := range validationError.Errors {

		if !strings.HasPrefix(basicError.Error, "missing properties:") && (basicError.InstanceLocation == "" || basicError.Error == "" || strings.HasPrefix(basicError.Error, "doesn't validate with")) {
			continue
		}
		if len(validationErrors) > 0 && validationErrors[len(validationErrors)-1].InstanceLocation == basicError.InstanceLocation {
			validationErrors[len(validationErrors)-1].Error += ", " + basicError.Error
		} else {
			failedValue := model.FindValue(originalObject, strings.Split(basicError.InstanceLocation, "/")[1:])
			_, mapOk := failedValue.(map[string]interface{})
			_, sliceOk := failedValue.([]interface{})
			if mapOk || sliceOk {
				failedValue = nil
			}
			// Create a ValidatorError from the jsonschema.BasicError
			validationError := validation.ValidatorError{
				KeywordLocation:         basicError.KeywordLocation,
				AbsoluteKeywordLocation: basicError.AbsoluteKeywordLocation,
				InstanceLocation:        basicError.InstanceLocation,
				Error:                   basicError.Error,
				FailedValue:             failedValue,
			}
			validationErrors = append(validationErrors, validationError)
		}
	}
	return validationErrors

}
