package validationstore_test

import (
	"context"
	"testing"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/pkg/common"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

const (
	validationPath           = "../../../test/e2e/scenarios/remote-validations/validation.opa.yaml"
	componentPath            = "../../../test/e2e/scenarios/remote-validations/component-definition.yaml"
	executableValidationPath = "../../../test/e2e/scenarios/create-resources/validation.yaml"
)

func generateValidation(t *testing.T, path string) common.Validation {
	validationBytes, err := common.ReadFileToBytes(path)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	var validation common.Validation
	err = validation.UnmarshalYaml(validationBytes)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	return validation
}

func TestNewValidationStore(t *testing.T) {
	v := validationstore.NewValidationStore()
	if v == nil {
		t.Error("Expected a new ValidationStore, but got nil")
	}
}

func TestNewValidationStoreFromBackMatter(t *testing.T) {
	backMatter := oscalTypes.BackMatter{}
	v := validationstore.NewValidationStoreFromBackMatter(backMatter)
	if v == nil {
		t.Error("Expected a new ValidationStore from back matter, but got nil")
	}
}

func TestAddValidation(t *testing.T) {
	validation := generateValidation(t, validationPath)
	v := validationstore.NewValidationStore()

	id, err := v.AddValidation(&validation)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if id == "" {
		t.Error("Expected a non-empty ID, but got an empty string")
	}
}

func TestGetLulaValidation(t *testing.T) {
	validation := generateValidation(t, validationPath)
	v := validationstore.NewValidationStore()
	id, _ := v.AddValidation(&validation)
	lulaValidation, err := v.GetLulaValidation(id)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if lulaValidation == nil {
		t.Error("Expected a LulaValidation, but got nil")
	}
}

func TestDryRun(t *testing.T) {
	validation := generateValidation(t, validationPath)
	executableValidation := generateValidation(t, executableValidationPath)

	tests := []struct {
		name               string
		validations        []common.Validation
		expectedExecutable bool
	}{
		{
			name: "No executable validations",
			validations: []common.Validation{
				validation,
				validation,
			},
			expectedExecutable: false,
		},
		{
			name: "Some executable validations",
			validations: []common.Validation{
				validation,
				executableValidation,
			},
			expectedExecutable: true,
		},
		{
			name: "All executable validations",
			validations: []common.Validation{
				executableValidation,
				executableValidation,
			},
			expectedExecutable: true,
		},
		{
			name:               "No validations",
			validations:        []common.Validation{},
			expectedExecutable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validationstore.NewValidationStore()
			for _, validation := range tt.validations {
				_, err := v.AddValidation(&validation)
				require.NoError(t, err)
			}

			executable, _ := v.DryRun()
			if executable != tt.expectedExecutable {
				t.Errorf("Expected %t, but got %t", tt.expectedExecutable, executable)
			}
		})
	}
}

func TestRunValidations(t *testing.T) {
	message.NoProgress = true
	validation := types.CreatePassingLulaValidation("sample-validation")

	tests := []struct {
		name                 string
		lulaValidations      []*types.LulaValidation
		expectedObservations int
	}{
		{
			name: "One validation",
			lulaValidations: []*types.LulaValidation{
				validation,
			},
			expectedObservations: 1,
		},
		{
			name: "Multiple validations",
			lulaValidations: []*types.LulaValidation{
				validation,
				validation,
			},
			expectedObservations: 2,
		},
		{
			name:                 "No validations",
			lulaValidations:      []*types.LulaValidation{},
			expectedObservations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validationstore.NewValidationStore()
			for _, validation := range tt.lulaValidations {
				v.AddLulaValidation(validation, uuid.NewUUID())
			}

			observations := v.RunValidations(context.Background(), true, false, "")
			if len(observations) != tt.expectedObservations {
				t.Errorf("Expected %d observations, but got %d", tt.expectedObservations, len(observations))
			}
		})
	}

	// Test that validations are stored as pointers -> store data across accesses
	t.Run("Validations store data", func(t *testing.T) {
		// Create a new validation store, add the validation
		validationUuid := uuid.NewUUID()
		v := validationstore.NewValidationStore()
		v.AddLulaValidation(validation, validationUuid)

		// Run the validations on the store
		_ = v.RunValidations(context.Background(), true, false, "")

		// Check that the validation data has been stored in v, state should be satisfied
		val, err := v.GetLulaValidation(validationUuid)
		require.NoError(t, err)
		require.NotNil(t, val)
		require.NotNil(t, val.Result)
		require.Equal(t, "satisfied", val.Result.State)
	})
}

func TestGetRelatedObservation(t *testing.T) {
	message.NoProgress = true
	validationPass := types.CreatePassingLulaValidation("passing-validation")
	validationFail := types.CreateFailingLulaValidation("failing-validation")
	v := validationstore.NewValidationStore()
	v.AddLulaValidation(validationPass, "1")
	v.AddLulaValidation(validationFail, "2")

	v.RunValidations(context.Background(), true, false, "")

	tests := []struct {
		name               string
		validationId       string
		expectedPass       bool
		relatedObservation bool
	}{
		{
			name:               "Pass Validation",
			validationId:       "1",
			expectedPass:       true,
			relatedObservation: true,
		},
		{
			name:               "Fail Validation",
			validationId:       "2",
			expectedPass:       false,
			relatedObservation: true,
		},
		{
			name:               "Non-existant Validation",
			validationId:       "3",
			expectedPass:       false,
			relatedObservation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relatedObservation, pass := v.GetRelatedObservation(tt.validationId)
			if pass != tt.expectedPass {
				t.Errorf("Expected %t, but got %t", tt.expectedPass, pass)
			}
			if tt.relatedObservation && relatedObservation.ObservationUuid == "" {
				t.Error("Expected observation UUID, but got empty")
			}
		})
	}
}

func TestRunTests(t *testing.T) {
	message.NoProgress = true
	ctx := context.Background()
	v := validationstore.NewValidationStore()

	validation := generateValidation(t, "./testdata/validation.yaml")
	validationWithTests := generateValidation(t, "./testdata/validation-with-tests.yaml")

	idValidation, err := v.AddValidation(&validation)
	require.NoError(t, err)
	idValidationWithTests, err := v.AddValidation(&validationWithTests)
	require.NoError(t, err)

	// Run validations to populate domain resources for tests
	v.RunValidations(ctx, true, false, "")

	// Run tests
	testReport := v.RunTests(ctx)

	// Validate the content of the test report
	reportValidation, ok := testReport[idValidation]
	require.True(t, ok)

	// no tests defined, should be empty
	assert.Equal(t, 0, len(reportValidation.TestResults))

	reportValidationWithTests, ok := testReport[idValidationWithTests]
	require.True(t, ok)

	// 2 tests defined, should have 2 results, both passing
	assert.Equal(t, 2, len(reportValidationWithTests.TestResults))
	assert.True(t, reportValidationWithTests.TestResults[0].Pass)
	assert.True(t, reportValidationWithTests.TestResults[1].Pass)
}
