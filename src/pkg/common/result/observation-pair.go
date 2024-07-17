package result

import (
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

type ObservationPair struct {
	stateChange         StateChange
	satisfied           bool
	name                string
	observation         string
	comparedObservation string
}

// NewObservationPair -> create a new observation pair from a given observation and slice of comparedObservations
func NewObservationPair(observation *oscalTypes_1_1_2.Observation, comparedObservations []*oscalTypes_1_1_2.Observation) *ObservationPair {
	// Calculate the state change
	var state StateChange
	var result bool
	var observationRemarks, comparedObservationRemarks, name string
	comparedObservation := findObservation(observation, comparedObservations)
	if observation == nil && comparedObservation == nil {
		state = UNCHANGED
	} else if observation != nil && comparedObservation == nil {
		state = NEW
		observationRemarks = getRemarks(observation.RelevantEvidence)
		name = strings.TrimPrefix(observation.Description, "[TEST]: ")
	} else if observation == nil && comparedObservation != nil {
		state = REMOVED
		comparedObservationRemarks = getRemarks(comparedObservation.RelevantEvidence)
		name = strings.TrimPrefix(comparedObservation.Description, "[TEST]: ")
	} else {
		// Check if the observation has changed
		observationRemarks = getRemarks(observation.RelevantEvidence)
		comparedObservationRemarks = getRemarks(comparedObservation.RelevantEvidence)
		name = strings.TrimPrefix(observation.Description, "[TEST]: ")
		state = getStateChange(observation, comparedObservation)
		result = getObservationResult(observation.RelevantEvidence)
	}

	return &ObservationPair{
		stateChange:         state,
		satisfied:           result,
		name:                name,
		observation:         observationRemarks,
		comparedObservation: comparedObservationRemarks,
	}
}

// findObservation finds an observation in a slice of observations
func findObservation(observation *oscalTypes_1_1_2.Observation, observations []*oscalTypes_1_1_2.Observation) *oscalTypes_1_1_2.Observation {
	for _, comparedObservation := range observations {
		if observation.Description == comparedObservation.Description {
			return comparedObservation
		}
	}
	return nil
}

// getStateChange compares the relevant evidence of two observations and calculates the state change between the two
func getStateChange(observation *oscalTypes_1_1_2.Observation, comparedObservation *oscalTypes_1_1_2.Observation) StateChange {
	var state StateChange = UNCHANGED
	relevantEvidence := observation.RelevantEvidence
	comparedRelevantEvidence := comparedObservation.RelevantEvidence

	if relevantEvidence == nil {
		if comparedRelevantEvidence != nil {
			state = REMOVED
		}
	} else {
		if comparedRelevantEvidence == nil {
			state = NEW
		} else {
			state = compareRelevantEvidence(relevantEvidence, comparedRelevantEvidence)
		}
	}

	return state
}

func compareRelevantEvidence(relevantEvidence *[]oscalTypes_1_1_2.RelevantEvidence, comparedRelevantEvidence *[]oscalTypes_1_1_2.RelevantEvidence) StateChange {
	var state StateChange = UNCHANGED

	reResults := getObservationResult(relevantEvidence)
	compReResults := getObservationResult(comparedRelevantEvidence)

	if reResults && !compReResults {
		state = NOT_SATISFIED_TO_SATISFIED
	} else if !reResults && compReResults {
		state = SATISFIED_TO_NOT_SATISFIED
	}

	return state
}

func getObservationResult(relevantEvidence *[]oscalTypes_1_1_2.RelevantEvidence) bool {
	var satisfied bool
	if relevantEvidence != nil {
		for _, re := range *relevantEvidence {
			if !strings.Contains(re.Description, "not-satisfied") {
				satisfied = true
			}
		}
	}
	return satisfied
}

func getRemarks(relevantEvidence *[]oscalTypes_1_1_2.RelevantEvidence) string {
	var remarks string
	if relevantEvidence != nil {
		remarks = (*relevantEvidence)[0].Remarks
	}
	return remarks
}
