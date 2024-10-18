package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/console"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

func TestConsoleCommand(t *testing.T) {
	message.NoProgress = true

	test := func(t *testing.T, args ...string) error {
		rootCmd := console.ConsoleCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	t.Run("Console run with no file args - error", func(t *testing.T) {
		err := test(t)
		require.Error(t, err)
	})

	t.Run("Console run with invalid oscal file - error", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/validation/validation.opa.yaml")
		require.Error(t, err)
	})

	t.Run("Console run with invalid output component file - error", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/oscal/valid-component.yaml", "-c", "../../unit/common/validation/validation.opa.yaml")
		require.Error(t, err)
	})

	t.Run("Console run with invalid output asessement results file - error", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/oscal/valid-component.yaml", "-a", "../../unit/common/validation/validation.opa.yaml")
		require.Error(t, err)
	})
}
