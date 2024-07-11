package validationResult_test

import (
	"errors"
	"testing"
	"time"

	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
	validationResult "github.com/defenseunicorns/lula/src/pkg/common/validation-result"
)

func TestNewNonSchemaValidationError(t *testing.T) {
	err := errors.New("test error")
	documentType := "testDocument"
	result := validationResult.NewNonSchemaValidationError(err, documentType)

	if result.Valid {
		t.Errorf("Expected Valid to be false, got true")
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Error != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", result.Errors[0].Error)
	}
	if result.Errors[0].AbsoluteKeywordLocation != validationResult.NON_SCHEMA_ERROR_ABSOLUTE_KEYWORD_LOCATION {
		t.Errorf("Expected AbsoluteKeywordLocation '%s', got '%s'", validationResult.NON_SCHEMA_ERROR_ABSOLUTE_KEYWORD_LOCATION, result.Errors[0].AbsoluteKeywordLocation)
	}
	if result.Metadata.DocumentType != documentType {
		t.Errorf("Expected DocumentType '%s', got '%s'", documentType, result.Metadata.DocumentType)
	}
	if time.Since(result.TimeStamp) > time.Second {
		t.Errorf("TimeStamp is too old")
	}
}

func TestIsNonSchemaValidationError(t *testing.T) {
	err := errors.New("test error")
	documentType := "testDocument"
	result := validationResult.NewNonSchemaValidationError(err, documentType)

	if !validationResult.IsNonSchemaValidationError(result) {
		t.Errorf("Expected IsNonSchemaValidationError to be true, got false")
	}

	// Test with a different error location
	result.Errors[0].AbsoluteKeywordLocation = "different-location"
	if validationResult.IsNonSchemaValidationError(result) {
		t.Errorf("Expected IsNonSchemaValidationError to be false, got true")
	}

	// Test with multiple errors
	result.Errors = append(result.Errors, oscalValidation.ValidatorError{Error: "another error"})
	if validationResult.IsNonSchemaValidationError(result) {
		t.Errorf("Expected IsNonSchemaValidationError to be false, got true")
	}
}

func TestGetNonSchemaError(t *testing.T) {
	err := errors.New("test error")
	documentType := "testDocument"
	result := validationResult.NewNonSchemaValidationError(err, documentType)

	extractedErr := validationResult.GetNonSchemaError(result)
	if extractedErr == nil {
		t.Errorf("Expected non-nil error, got nil")
	}
	if extractedErr.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", extractedErr.Error())
	}

	// Test with a different error location
	result.Errors[0].AbsoluteKeywordLocation = "different-location"
	extractedErr = validationResult.GetNonSchemaError(result)
	if extractedErr != nil {
		t.Errorf("Expected nil error, got '%s'", extractedErr.Error())
	}

	// Test with multiple errors
	result.Errors = append(result.Errors, oscalValidation.ValidatorError{Error: "another error"})
	extractedErr = validationResult.GetNonSchemaError(result)
	if extractedErr != nil {
		t.Errorf("Expected nil error, got '%s'", extractedErr.Error())
	}
}
