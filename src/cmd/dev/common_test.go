package dev_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
)

func TestGetObservationByUuid(t *testing.T) {
	t.Parallel()

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-assessment-results-with-resources.yaml")
	assessment := oscalModel.AssessmentResults
	require.NotNil(t, assessment)

	t.Run("Test get observation by uuid - found", func(t *testing.T) {
		observation, err := dev.GetObservationByUuid(assessment, "92cb3cad-bbcd-431a-aaa9-cd47275a3982")
		require.NoError(t, err)
		require.NotNil(t, observation)
	})

	t.Run("Test get observation by uuid - not found", func(t *testing.T) {
		observation, err := dev.GetObservationByUuid(assessment, "invalid-uuid")
		assert.Nil(t, observation)
		require.ErrorContains(t, err, "observation with uuid invalid-uuid not found")
	})

}
