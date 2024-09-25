package test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/stretchr/testify/assert"
)

const (
	validMultiComponentPath = "../unit/common/oscal/valid-multi-component.yaml"
	catalogPath             = "../unit/common/oscal/catalog.yaml"
)

// TestLulaReportValidComponent checks that the 'lula report' command works with a valid component definition.
func TestLulaReportValidComponent(t *testing.T) {
	// Disable progress indicators and other extra formatting
	message.NoProgress = true

	var outBuf, errBuf bytes.Buffer

	cmd := exec.Command("lula", "report", "-f", validMultiComponentPath, "--file-format", "table")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	// Check for errors in command execution.
	assert.NoError(t, err, "Expected no error from `lula report` with valid component definition")
}

// TestLulaReportCatalog checks that the 'lula report' command fails gracefully with a catalog file.
// OSCAL Catalogs are not currently supported by 'lula report' command yet.
func TestLulaReportCatalog(t *testing.T) {
	// Disable progress indicators and other extra formatting
	message.NoProgress = true

	var outBuf, errBuf bytes.Buffer

	cmd := exec.Command("lula", "report", "-f", catalogPath, "--file-format", "table")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	assert.Error(t, err, "error running report")
}
