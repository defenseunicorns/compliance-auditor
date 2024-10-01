package oscal_test

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestGetType(t *testing.T) {
	test := func(t *testing.T, model oscal.Profile, expected string) {
		t.Helper()

		got := model.GetType()

		if got != expected {
			t.Fatalf("Expected %s - got %s\n", expected, got)
		}
	}

	t.Run("Test populated model", func(t *testing.T) {

		var profile = oscal.Profile{
			Model: &oscalTypes.Profile{},
		}

		test(t, profile, "profile")
	})

	t.Run("Test unpopulated model", func(t *testing.T) {

		var profile = oscal.Profile{}

		test(t, profile, "profile")
	})
}
