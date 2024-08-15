package schemas

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/model"
	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
)

//go:embed *.json
var Schemas embed.FS

const (
	SCHEMA_SUFFIX      = ".json"
	SCHEMA_PATH_PREFIX = "https://github.com/defenseunicorns/lula/tree/main/src/pkg/schemas/"
)

func PrefixSchema(path string) string {
	if !strings.HasSuffix(path, SCHEMA_SUFFIX) {
		path = path + SCHEMA_SUFFIX
	}
	return path
}

func GetSchemaPath(schemaName string) string {
	return SCHEMA_PATH_PREFIX + schemaName + SCHEMA_SUFFIX
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

func Validate(schema string, data model.InterfaceOrBytes) oscalValidation.ValidationResult {
	validationParams := &oscalValidation.ValidationParams{
		ModelType: schema,
	}

	schemaBytes, err := GetSchema(schema)
	if err != nil {
		return *oscalValidation.NewNonSchemaValidationError(err, validationParams)
	}
	schemaData, err := model.CoerceToJsonMap(schemaBytes)
	if err != nil {
		return *oscalValidation.NewNonSchemaValidationError(err, validationParams)
	}

	modelData, err := model.CoerceToJsonMap(data)
	if err != nil {
		return *oscalValidation.NewNonSchemaValidationError(err, validationParams)
	}

	validationParams.SchemaData = schemaData
	validationParams.SchemaPath = GetSchemaPath(schema)
	validationParams.ModelData = modelData

	result, _ := oscalValidation.ValidateFromParams(validationParams)

	return *result
}
