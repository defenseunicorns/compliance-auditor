# Console

The Lula Console is a text-based terminal user interface that allows users to interact with the OSCAL documents in a more intuitive and visual way.

Currently, only the **Component Definition** and **Assessment Results** models are supported in the Console.

 * See the sub-pages for more information on interacting with each specific OSCAL model in the Console.

>[!NOTE]
>The Console is currently in development and views are subject to change.

## Usage

To open the Console with particular OSCAL models:
```shell
lula console -f /path/to/oscal-component.yaml,/path/to/oscal-component-2.yaml,/path/to/assessment-results.yaml
```
The `-f` (or `--input-files`) flag can be used to specify multiple OSCAL model file paths to be loaded into the Console.

### Writing to Output

#### Component Definition

The Console supports writing updated `remarks` and `description` fields to the `component-definition` model.

To include an output file to save any changes made to the component definition, use the `--component-output` or `-c`flag:
```shell
lula console -f /path/to/oscal-component.yaml -c /path/to/output.yaml
```

If no output file is specified and a single component definition is passed, the provided component definition will be overwritten. If multiple component definitions are passed and no output file is specified, the Console will default to `component.yaml` in the current working directory.

#### Assessment Results

The Console supports writing assessment results that are generated from the Component Definition `validate` execution.

To specify an output file to write assessment results, include `--assessment-output` or `-a`. If you provide no `-a`, but include a single assessment results file as inputs (via `--input-files`/`-f`), that will be used as the output file. If no assessment results file is provided (or multiple are provided), the Console will default to `assessment-results.yaml` in the current working directory.

Example of running console including the optionality to write assessment results to the provided file:
```shell
lula console -f /path/to/oscal-component.yaml -a /path/to/output.yaml
```

## Keys

The Console responds to the following keys for navigation and interaction (each sub-model has additional key response, see respective help views for more information):

| Key | Description |
|-----|-------------|
| `?` | Toggle help |
| `ctrl+c` | Quit |
| `tab` | Tab right between models |
| `shift+tab` | Tab left between models |