package oscal

import (
	"os"
	"testing"
)

func TestMultiValidate(t *testing.T) {
	multiModelData, err := os.ReadFile("../../../test/unit/common/oscal/multi-model.yaml")
	if err != nil {
		t.Fatalf("os.ReadFile failed: %v", err)
	}
	t.Run("Test Multi Validate", func(t *testing.T) {
		data := multiModelData
		err := multiModelValidate(data)
		if err != nil {
			t.Fatalf("multiModelValidate failed: %v", err)
		}
	})
}
