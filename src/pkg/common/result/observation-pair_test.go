package result_test

import (
	"fmt"
	"testing"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/result"
)

func createObservation(description, satisfaction string) *oscalTypes_1_1_2.Observation {
	return &oscalTypes_1_1_2.Observation{
		UUID:        uuid.NewUUID(),
		Description: description,
		RelevantEvidence: &[]oscalTypes_1_1_2.RelevantEvidence{
			{
				Description: fmt.Sprintf("Result: %s", satisfaction),
				Remarks:     "Some remarks about this observation",
			},
		},
	}
}

func TestCreateObservationPairs(t *testing.T) {
	// tests different variations of observation pairs
	tests := []struct {
		name                string
		observations        []*oscalTypes_1_1_2.Observation
		compareObservations []*oscalTypes_1_1_2.Observation
		expectedPairs       int
		expectedStateChange map[string]result.StateChange
	}{
		{
			name: "One observation pair, not satisfied to satisfied",
			observations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "satisfied"),
			},
			compareObservations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "not-satisfied"),
			},
			expectedPairs: 1,
			expectedStateChange: map[string]result.StateChange{
				"test-1": result.NOT_SATISFIED_TO_SATISFIED,
			},
		},
		{
			name: "One observation pair, satisfied to not-satisfied",
			observations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "not-satisfied"),
			},
			compareObservations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "satisfied"),
			},
			expectedPairs: 1,
			expectedStateChange: map[string]result.StateChange{
				"test-1": result.SATISFIED_TO_NOT_SATISFIED,
			},
		},
		{
			name: "Two observation pairs",
			observations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "satisfied"),
			},
			compareObservations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "satisfied"),
				createObservation("test-2", "not-satisfied"),
			},
			expectedPairs: 2,
			expectedStateChange: map[string]result.StateChange{
				"test-1": result.UNCHANGED,
				"test-2": result.REMOVED,
			},
		},
		{
			name: "Three observation pairs",
			observations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-1", "satisfied"),
				createObservation("test-3", "not-satisfied"),
			},
			compareObservations: []*oscalTypes_1_1_2.Observation{
				createObservation("test-2", "not-satisfied"),
				createObservation("test-3", "not-satisfied"),
			},
			expectedPairs: 3,
			expectedStateChange: map[string]result.StateChange{
				"test-1": result.NEW,
				"test-2": result.REMOVED,
				"test-3": result.UNCHANGED,
			},
		},
		{
			name:                "No observation pairs",
			observations:        []*oscalTypes_1_1_2.Observation{},
			compareObservations: []*oscalTypes_1_1_2.Observation{},
			expectedPairs:       0,
			expectedStateChange: map[string]result.StateChange{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observationPairs := result.CreateObservationPairs(tt.observations, tt.compareObservations)

			if len(observationPairs) != tt.expectedPairs {
				t.Errorf("Expected %d pairs, but got %d pairs", tt.expectedPairs, len(observationPairs))
			}

			for _, op := range observationPairs {
				if op.StateChange != tt.expectedStateChange[op.Name] {
					t.Errorf("Expected %s, but got %s", tt.expectedStateChange[op.Name], op.StateChange)
				}
			}
		})
	}
}
