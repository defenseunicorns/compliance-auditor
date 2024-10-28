package result

import (
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

type ObservationPair struct {
	StateChange             StateChange
	Satisfied               bool
	Name                    string
	Observation             string
	ObservationUuid         string
	ComparedObservation     string
	ComparedObservationUuid string
}

// CreateObservationPairs creates a slice of observation pairs from a slice of observations and compared observations
func CreateObservationPairs(observations []*oscalTypes_1_1_2.Observation, comparedObservations []*oscalTypes_1_1_2.Observation) []*ObservationPair {
	observationPairs := make([]*ObservationPair, 0)

	// Add all observations to the observation pairs
	for _, observation := range observations {
		comparedObservation := findObservation(observation, comparedObservations)
		observationPair := newObservationPair(observation, comparedObservation)
		observationPairs = append(observationPairs, observationPair)
	}

	// Add all compared observations that are not in the observations
	for _, comparedObservation := range comparedObservations {
		observation := findObservation(comparedObservation, observations)
		if observation == nil {
			observationPair := newObservationPair(nil, comparedObservation)
			observationPairs = append(observationPairs, observationPair)
		}
	}
	return observationPairs
}

// NewObservationPair -> create a new observation pair from a given observation and slice of comparedObservations
func newObservationPair(observation *oscalTypes_1_1_2.Observation, comparedObservation *oscalTypes_1_1_2.Observation) *ObservationPair {
	// Calculate the state change
	var state StateChange
	var result bool
	var observationRemarks, comparedObservationRemarks, observationUuid, comparedObservationUuid, name string
	prefix := "[TEST]: "

	if observation != nil {
		observationUuid = observation.UUID
		name = strings.TrimPrefix(observation.Description, prefix)
		observationRemarks = getRemarks(observation.RelevantEvidence)
		result = getObservationResult(observation.RelevantEvidence)
		if comparedObservation == nil {
			state = NEW
		} else {
			comparedObservationUuid = comparedObservation.UUID
			comparedObservationRemarks = getRemarks(comparedObservation.RelevantEvidence)
			state = getStateChange(observation, comparedObservation)
		}
	} else {
		if comparedObservation != nil {
			comparedObservationUuid = comparedObservation.UUID
			name = strings.TrimPrefix(comparedObservation.Description, prefix)
			comparedObservationRemarks = getRemarks(comparedObservation.RelevantEvidence)
			state = REMOVED
		} else {
			state = UNCHANGED
		}
	}

	return &ObservationPair{
		StateChange:             state,
		Satisfied:               result,
		Name:                    name,
		Observation:             observationRemarks,
		ObservationUuid:         observationUuid,
		ComparedObservation:     comparedObservationRemarks,
		ComparedObservationUuid: comparedObservationUuid,
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
