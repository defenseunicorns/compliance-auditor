package validationResult

import (
	"errors"
	"time"

	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
)

// NON_SCHEMA_ERROR_ABSOLUTE_KEYWORD_LOCATION is the absolute keyword location for non-schema errors
const NON_SCHEMA_ERROR_ABSOLUTE_KEYWORD_LOCATION = "non-schema-error"

// NewNonSchemaValidationError creates a system validation error
func NewNonSchemaValidationError(err error, documentType string) oscalValidation.ValidationResult {
	return oscalValidation.ValidationResult{
		Valid:     false,
		TimeStamp: time.Now(),
		Errors: []oscalValidation.ValidatorError{
			{
				Error:                   err.Error(),
				AbsoluteKeywordLocation: NON_SCHEMA_ERROR_ABSOLUTE_KEYWORD_LOCATION,
			},
		},
		Metadata: oscalValidation.ValidationResultMetadata{
			DocumentType: documentType,
		},
	}
}

// IsNonSchemaValidationError checks if the result is a system validation error
func IsNonSchemaValidationError(result oscalValidation.ValidationResult) bool {
	return len(result.Errors) == 1 && result.Errors[0].AbsoluteKeywordLocation == NON_SCHEMA_ERROR_ABSOLUTE_KEYWORD_LOCATION
}

// GetNonSchemaError extracts the system validation error
// If the result is not a system validation error or if there are no errors, return nil
func GetNonSchemaError(result oscalValidation.ValidationResult) error {
	if !IsNonSchemaValidationError(result) {
		return nil
	}
	return errors.New(result.Errors[0].Error)
}
