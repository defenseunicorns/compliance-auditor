# Component Definition

The Component Definition view currently allows for a read and limited write experience of the OSCAL Component Definition model. The view is tailored to the usage of the Component Definition in the context of Lula, and is not intended to be a comprehensive view of the model.

## Usage

To view an OSCAL Component Definition model in the Console:
```shell
lula console -f /path/to/oscal-component.yaml
```
The `oscal-component.yaml` will need to be a valid OSCAL model - to use with the Component Definition view, it must contain the `component-definition` top level key.

To include an output file to save any changes made to the component definition, use the `--component-output` or `-c`flag:
```shell
lula console -f /path/to/oscal-component.yaml -c /path/to/output.yaml
```

> [!Note] 
> Several component definition models can be passed into the console, via `-f` in a comma-separated list. For multiple component definitions, the output file will default to `component.yaml` unless specified.

## Keys

The Component Definition model responds to the following keys for navigation and interaction (some widgets have additional key response, see respective help views for more information):

| Key | Description |
|-----|-------------|
| `?` | Toggle help |
| `ctrl+c` | Quit |
| `tab` | Tab right between models |
| `shift+tab` | Tab left between models |
| `←/h` | Navigate left in model|
| `→/l` | Navigate right in model |
| `↑/k` | Move up in list OR scroll up in panel |
| `↓/j` | Move down in list OR scroll up in panel |
| `/` | Filter list |
| `↳` | Select item |
| `e` | Edit available fields (remarks and description) |
| `ctrl+s` | Save changes (Note: this may overwrite the original file if an output file unspecified) |
| `esc` | Cancel |

During console viewing, the top-right corner will display the help keys availble in the context of the selected widget. When an overlay is open, the help keys will be displayed in the overlay.

## Views

### Read-Only Navigation

<img align="right" src="../../images/component-defn-console-read.gif" alt="component definition console read" style="width:100%; height:auto;">

### Editing Remarks and Description

<img align="right" src="../../images/component-defn-console-edit.gif" alt="component definition console edit" style="width:100%; height:auto;">
