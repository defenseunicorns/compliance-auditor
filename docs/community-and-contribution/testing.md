# Testing

## Execution

There are multiple layers of tests that can be run in isolation or collectively as a whole. 

- `make test` will execute the full suite of tests
  - This requires an environment with `kind` available
- `make test-e2e` will execute the end-to-end tests for CLI and Kubernetes testing
  - This requires an environment with `kind` available
- `make test-cmd` tests the CLI tests
  - Does not require additional infrastructure
- `make test-unit` runs the unit tests
  - Does not require additional infrastructure

## Test Data

Testing artifacts are stored in relation to the tests being run or centralized for access across multiple tests. All `.golden` files are generated with `go test <path> -update` and should not be modified manually. 

### Console Testing

The Lula Console is a text-based terminal user interface that allows users to interact with the OSCAL documents and is written using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) library. 

To test the Lula Console, we've implemented [teatest](https://pkg.go.dev/github.com/charmbracelet/x/exp/teatest), which allows us to generate "golden" snapshots of the console output, then ensure the test results match that expected output.

#### Usage

To update the golden snapshot for the Lula Console, run the following command:

```shell
go test ./src/internal/tui/model_test.go -update 
```

This will update the golden snapshot files in the `testdata` directory.