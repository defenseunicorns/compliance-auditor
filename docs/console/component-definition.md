# Component Definition

The Component Definition view currently allows for a read and limited write experience of the OSCAL Component Definition model. The view is tailored to the usage of the Component Definition in the context of Lula, and is not intended to be a comprehensive view of the model.

## Usage

The Component Definition model supports the following views:
 * [Read-Only Navigation](./component-definition.md#read-only-navigation)
 * [Editing Remarks and Description](./component-definition.md#editing-remarks-and-description)
 * [Validation Detail](./component-definition.md#validation-detail)

## Keys

The Component Definition model responds to the following keys for navigation and interaction (some widgets have additional key response, see respective help views for more information):

| Key | Description |
|-----|-------------|
| `?` | Toggle help |
| `ctrl+c` | Quit |
| `tab` | Tab right between models |
| `shift+tab` | Tab left between models |
| `←/h` | Navigate left across widgets model|
| `→/l` | Navigate right across widgets model |
| `↑/k` | Move up in list OR scroll up in panel |
| `↓/j` | Move down in list OR scroll up in panel |
| `/` | Filter list |
| `↳` | Select item |
| `e` | Edit available fields (remarks and description) |
| `d` | Detail available fields (validations) |
| `ctrl+s` | Save changes (Note: this may overwrite the original file if an output file unspecified) |
| `esc` | Cancel |

During console viewing, the top-right corner will display the help keys availble in the context of the selected widget. When an overlay is open, the help keys will be displayed in the overlay.

## Views

### Read-Only Navigation

The model can be sorted by Component, Framework, and Control. Additional data/features provided by the Component Definition OSCAL Model are not currently supported for viewing in the Console.

<img align="right" src="../../images/component-defn-console-read.gif" alt="component definition console read" style="width:100%; height:auto;">

### Editing Remarks and Description

Limited editing of the remarks and description is supported. Once changes are made, to be persisted back to the file, the data will need to be saved via the `ctrl+s` key.

<img align="right" src="../../images/component-defn-console-edit.gif" alt="component definition console edit" style="width:100%; height:auto;">

### Validation Detail

The Validation Detail is a view that displays a somewhat curated version of the Lula Validation. It is intended to be a quick view of the validation, and is not a one-to-one representation.

<img align="right" src="../../images/component-defn-console-validation-detail" alt="component definition console validation detail" style="width:100%; height:auto;">
