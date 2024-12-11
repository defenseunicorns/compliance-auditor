package dev

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzRunSingleValidation(f *testing.F) {
	// including a fully populated valid-component.yaml to exercise a fuzzy version of passing valid-but-wrong-for-this-command oscal data.
	for _, tc := range []string{"../../test/unit/common/oscal/valid-component.yaml", "../../test/unit/common/validation/multi.validation.yaml"} {
		bytes, err := os.ReadFile(tc)
		require.NoError(f, err)

		f.Add(bytes)
	}

	f.Fuzz(func(t *testing.T, a []byte) {
		RunSingleValidation(context.Background(), a)
	})
}

func FuzzDevValidate(f *testing.F) {
	bytes, err := os.ReadFile("../../test/unit/common/validation/multi.validation.yaml")
	require.NoError(f, err)

	drs, err := os.ReadFile("../../test/unit/common/resources/valid-resources.json")
	require.NoError(f, err)
	f.Add(bytes, drs)

	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		DevValidate(context.Background(), a, b, false, nil)
	})
}
