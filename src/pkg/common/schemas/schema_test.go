package schemas_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common/schemas"
	validationResult "github.com/defenseunicorns/lula/src/pkg/common/validation-result"
)

func TestToMap(t *testing.T) {
	t.Parallel() // Enable parallel execution of tests
	t.Run("Should return a map with all the schemas", func(t *testing.T) {
		t.Parallel() // Enable parallel execution of subtests
		fileMap, err := schemas.ToMap()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(fileMap) == 0 {
			t.Errorf("Expected map to have some files, got %d", len(fileMap))
		}

		// Check if all files have the correct JSON schema suffix and are not directories
		for fileName, fileInfo := range fileMap {
			if fileInfo.IsDir() {
				t.Errorf("Expected file but got directory for %s", fileName)
			}
			if fs.ValidPath(fileName) && !schemas.HasSchema(fileName) {
				t.Errorf("File %s does not have the correct schema suffix", fileName)
			}
		}
	})
}

func TestHasSchema(t *testing.T) {
	t.Parallel() // Enable parallel execution of tests
	t.Run("Should detect schema suffix correctly", func(t *testing.T) {
		t.Parallel() // Enable parallel execution of subtests
		validSchema := "validation.json"
		invalidSchema := "validation.txt"

		if !schemas.HasSchema(validSchema) {
			t.Errorf("Expected true for %s, got false", validSchema)
		}
		if schemas.HasSchema(invalidSchema) {
			t.Errorf("Expected false for %s, got true", invalidSchema)
		}
	})
}

func TestListSchemas(t *testing.T) {
	t.Parallel() // Enable parallel execution of tests
	t.Run("Should list all schemas", func(t *testing.T) {
		t.Parallel() // Enable parallel execution of subtests
		schemasList, err := schemas.ListSchemas()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(schemasList) == 0 {
			t.Errorf("Expected non-empty schema list, got empty")
		}
	})
}

func TestValidate(t *testing.T) {
	t.Parallel() // Enable parallel execution of tests
	validationPath := "../../../test/unit/common/validation/opa.validation.yaml"
	validationData, err := os.ReadFile(validationPath)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Run("Should validate a schema", func(t *testing.T) {
		t.Parallel() // Enable parallel execution of subtests
		schema := "validation"
		result := schemas.Validate(schema, validationData)
		if validationResult.GetNonSchemaError(result) != nil {
			t.Errorf("expected result to be valid, got %v", result)
		}
	})

	t.Run("Should return an ValidationResult if the schema is missing required properties", func(t *testing.T) {
		t.Parallel() // Enable parallel execution of subtests
		schema := "validation"
		result := schemas.Validate(schema, []byte("{\n\t\"name\": \"test\"\n}"))
		if result.Valid == true {
			t.Errorf("expected result to be invalid, got %v", result)
		}
	})
}
