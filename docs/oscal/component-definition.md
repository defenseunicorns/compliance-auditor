# Component Definition

A [Component Definition](https://pages.nist.gov/OSCAL/resources/concepts/layer/implementation/component-definition/) is an OSCAL model for capturing control information that pertains to a specific component/capability of a potential system. It can largely be considered the modular and re-usable model for use across many systems. In Lula, the `validate` command will process a `component-definition`, iterate through all `implemented-requirements` to discover Lula validations, and execute those validations to produce `observations`. 

## Components/Capabilities and Control-Implementations

The modularity of `component-definitions` allows for the specification of one to many components or capabilities that include one to many `control-implementations`.

By allowing for many `control-implementations`, a given component or capability can have information as to its compliance with many different regulatory standards. 

## Structure
The primary structure for Lula production and operations of `component-definitions` for determinism is as follows:
- Components/Capabilities are sorted by `title` in ascending order (Case Sensitive Sorting).
- Control Implementations are sorted by `source` in ascending order.
- Implemented Requirements are sorted by `control-id` in ascending order.
- Back Matter Resources are sorted by `title` in ascending order (Case Sensitive Sorting).

## Generation

Lula can generate OSCAL templates to help with authoring OSCAL artifacts. These generation processes help with initial build and maintenance of OSCAL artifacts and keeps with the vision of modular compliance artifacts that live with the source code. 

To generate a component definition, you need the following context:
- The catalog source `-c` or `--catalog-source`; IE `https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json`
- The controls from the catalog to map to `implemented-requirements` / `-r` or `--requirements`; `ac-1,ac-2,au-5`

The following command will generate a component definition with the above context:
```
lula generate component -c https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json -r ac-1,ac-2,au-5
```

There are optional flags that can be added to the command to generate a component definition - see [lula generate component](../cli-commands/lula_generate_component.md) for details. 

### Reproducibility  

The `lula generate` commands are meant to be reproducible and will auto-merge models based on filename. The intent for this generation is to make it easy to update a given model with automation and only inject human intervention as needed. An artifact generated with `lula generate` can be merged with a pre-existing artifact of the same model type. 

For component-definitions, see each individual `control-implementation` props for the `generation` prop. It should look like the following:
```yaml
props:
  - name: generation
    ns: https://docs.lula.dev/oscal/ns
    value: lula generate component --catalog-source https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json --component 'Component Title' --requirements ac-1,ac-3,ac-3.2,ac-4 --remarks assessment-objective
```

This `value` should mirror any required inputs in order to reproduce a given control implementation in a component. 

### Existing Data

The ability to retain data that is put into OSCAL artifacts is of utmost importance to this generation process and also a large feature of continued maintenance of these artifacts. Lula supports the ability to merge newly generated component definition templates into existing component definitions automatically. 

By specifying `--output` or `-o` and providing an existing file - Lula will perform a merge operation that only overwrites specific fields owned by automation.

Lula performs a match on the component title and the provided catalog source to determine placement and merge of the new implemented requirements. This can be used to updated exiting items or as a method to generation of a single artifacts that contains the data for many components or many control implementations. 

## Example 

```bash
./bin/lula generate component -c https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json -r ac-1,ac-3,ac-3.2,ac-4 -o oscal-component.yaml --remarks assessment-objective 
```