package testhelpers

import (
	"fmt"
	"os"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func OscalFromPath(t *testing.T, path string) *oscalTypes_1_1_2.OscalCompleteSchema {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	oscalModel, err := oscal.NewOscalModel(data)
	if err != nil {
		t.Fatalf("error creating oscal model from file: %v", err)
	}

	return oscalModel
}

func CreateTempFile(t *testing.T, ext string) *os.File {
	t.Helper()
	tempFile, err := os.CreateTemp("", fmt.Sprintf("tmp-*.%s", ext))
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	return tempFile
}
