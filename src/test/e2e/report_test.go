package test

import (
    "testing"

    "github.com/defenseunicorns/lula/src/cmd/report"
)

func TestReportCommand(t *testing.T) {

	// Define the test cases
	testCases := []struct {
		name       string
		inputFile  string
		fileFormat string
        expectErr  bool
	}{
		{
			name:       "Component Definition with one source and no framework props",
			inputFile:  "../../test/e2e/scenarios/report-data/valid-component.yaml",
			expectErr:  false,
		},
		{
			name:       "Component Definition with two sources and two framework props",
			inputFile:  "../../test/e2e/scenarios/report-data/valid-multi-component.yaml",
			expectErr:  false,
		},
        {
			name:       "Component Definition with two sources and two framework props with yaml file format",
			inputFile:  "../../test/e2e/scenarios/report-data/valid-multi-component.yaml",
			expectErr:  false,
            fileFormat: "yaml",
		},
        {
			name:       "Component Definition with two sources and two framework props with json file format",
			inputFile:  "../../test/e2e/scenarios/report-data/valid-multi-component.yaml",
			expectErr:  false,
            fileFormat: "json",
		},
		{
			name:       "Component Definition with one sources and one framework props",
			inputFile:  "../../test/e2e/scenarios/report-data/valid-component-with-framework.yaml",
			expectErr:  false,
		},
        {
            name:       "Catalog OSCAL Model",
            inputFile:  "../../test/e2e/scenarios/report-data/catalog.yaml",
            expectErr:  true,
        },
	}

    for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                // Call GenerateReport with the input file and format
                err := report.GenerateReport(tc.inputFile, tc.fileFormat)

                // Check if the result matches the expected outcome
                if tc.expectErr {
                    if err == nil {
                        t.Errorf("expected an error but got none for test case: %s", tc.name)
                    }
                } else {
                    if err != nil {
                        t.Errorf("did not expect an error but got one for test case: %s, error: %v", tc.name, err)
                    }
                }
            })
        }

}
