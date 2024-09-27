package test

import (
	"bytes"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/stretchr/testify/assert"
)

// TestLulaReportValidComponent checks that the 'lula report' command works with a valid component definition with multiple components and framework prop.
func TestLulaReportValidComponent(t *testing.T) {
	// Disable progress indicators and other extra formatting
	message.NoProgress = true

	// Setup the root command and buffers for capturing output
	rootCmd := cmd.RootCommand()
	rootCmd.SetArgs([]string{"report", "-f", "../unit/common/oscal/valid-multi-component.yaml", "--file-format", "table"})

	var outBuf, errBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&errBuf)

	// Execute the command
	err := rootCmd.Execute()

	// Check for errors in command execution.
	assert.NoError(t, err, "Expected no error from `lula report` with valid component definition")}

// TestLulaReportCatalog checks that the 'lula report' command fails gracefully with a catalog file since it is not currently supported.
// TODO: This test will need to be changed as more models are supported for reporting.
func TestLulaReportCatalog(t *testing.T) {
	// Disable progress indicators and other extra formatting
	message.NoProgress = true

	// Setup the root command and buffers for capturing output
	rootCmd := cmd.RootCommand()
	rootCmd.SetArgs([]string{"report", "-f", "../unit/common/oscal/catalog.yaml", "--file-format", "table"})

	var outBuf, errBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&errBuf)

	// Execute the command
	err := rootCmd.Execute()

	assert.Error(t, err, "Expected an error running report")}
