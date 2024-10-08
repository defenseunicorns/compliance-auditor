package tui_test

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/muesli/termenv"
)

const (
	timeout    = time.Second * 20
	maxRetries = 3
	height     = common.DefaultHeight
	width      = common.DefaultWidth
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

// TestNewOSCALModel tests that the NewOSCALModel creates the expected model from component definition file
func TestNewOSCALModel(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	model := tui.NewOSCALModel(oscalModel, "", tempLog)

	msgs := []tea.Msg{}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	if err != nil {
		t.Fatal(err)
	}
}
