# OSCAL Generation

Lula has the potential to provide a codified base for how to generate and maintain OSCAL through automation. This means that with a foundation built - Lula can continue to iterate on the methods for mapping & maintaining data that aligns with the intent of OSCAL and the standards/benchmarks involved. 

## Generic Generation Concepts

The generation process for OSCAL artifacts created and maintained by Lula should include the following:
- Specification of fields only maintained by automation (Challenge this one)
  - `implemented-requirements.remarks`
- Ability to maintain data that is added through manual interaction
  - This involves merging newly generated data with existing data
  - Involves two scenarios
    - No data exists - create a new document
    - Data exists - perform a merge of the new component with the existing data

## Scenarios for Manipulation of Component Definitions

### Component Definition Generation

This is the execution that involves creation of a new component-definition and then determining if there is an existing component-definition with which to merge the newly generated data.

This command will focus solely on a single component-definition `component` containing a single `control-implementation` - currently preventing further complexity. 

Current TODO:
- Ability to retain data in an existing OutputFile on re-generation
  - Easier said than done - check if the control-implementation exists
  - then find the delta of the controls
- Ability to detect an OSCAL manifest file (IE InputFile flag)
- wildcard match on requirements?

### Component Definition Import-Component-Definitions Compose

This involves importing other component-definition files from some location and performing a merge activity at all layers:
- Component
- Control Implementation
- Implemented Requirement
- Validation Link?

### Overlap
If we perform a merge of an original component-definition with a single component/control-implementation for Generation - then this could be re-used for import where merging could loop through this function.
What does this function really need? currently hardcodes the first component and the first control-implementation.
- In order to support merge with many to many of the above layers... we need the component and control-implementation to be passed in as a parameter?

## Example 

```bash
./bin/lula generate component -c https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json -r ac-1,ac-3,ac-3.2,ac-4 -o oscal-component.yaml --remarks assessment-objective -l debug
```