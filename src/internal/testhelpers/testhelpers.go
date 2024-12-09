package testhelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func OscalFromPath(t *testing.T, path string) *oscalTypes.OscalCompleteSchema {
	t.Helper()
	path = filepath.Clean(path)
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

// RunTestModelView runs a test model view with a given model and messages, impelements a retry loop if final model is nil
func RunTestModelView(t *testing.T, m tea.Model, reset func() tea.Model, msgs []tea.Msg, timeout time.Duration, maxRetries, height, width int) error {

	testModelView := func(t *testing.T, try int) (bool, error) {
		tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(width, height))

		for _, msg := range msgs {
			tm.Send(msg)
			time.Sleep(time.Millisecond * time.Duration(50*try))
		}

		if err := tm.Quit(); err != nil {
			return false, err
		}

		fm := tm.FinalModel(t, teatest.WithFinalTimeout(timeout))

		if fm == nil {
			return true, nil
		}

		teatest.RequireEqualOutput(t, []byte(fm.View()))

		return false, nil
	}

	for i := 0; i < maxRetries; i++ {
		retry, err := testModelView(t, i+1)
		if retry {
			if reset != nil {
				m = reset()
			}
			continue
		}
		if err != nil {
			return err
		}
		break
	}
	return nil
}
