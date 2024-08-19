# Report

The Lula `report` command will ingest OSCAL Component Definition models and will display the number of controls mapped on a per source and per framework. This will allow for uses to quickly see how many controls are mapped per source in the event of multiple `control-implementations`. This also allows for the Lula specific `prop` of `framework` which allows for a more custom mapping.

Use the command as follows:

For the default report

```bash
lula report -f oscal-component-definition.yaml
```

To change the display format to json or yaml

```bash
lula report -f oscal-component-definition.yaml --file-format json
```
