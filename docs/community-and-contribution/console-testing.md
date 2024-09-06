# Console Testing

The Lula Console is a text-based terminal user interface that allows users to interact with the OSCAL documents and is written using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) library. 

To test the Lula Console, we've implemented [teatest](https://pkg.go.dev/github.com/charmbracelet/x/exp/teatest), which allows us to generate "golden" snapshots of the console output, then ensure the test results match that expected output.

## Usage

To update the golden snapshot for the Lula Console, run the following command:

```shell
go test ./src/internal/tui/model_test.go -update 
```

This will update the golden snapshot files in the `testdata` directory.