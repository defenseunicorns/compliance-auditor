# OSCAL Generation

Lula has the potential to provide a codified base for how to generate and maintain OSCAL through automation. This means that with a foundation built - Lula can continue to iterate on the methods for mapping maintaining data the aligns with the intent of OSCAL and the standards/benchmarks involved. 

## Generic Generation Concepts

The generation process for OSCAL artifacts created and maintained by Lula should include the following:
- Specification of fields only maintained by automation (Challenge this one)
- Ability to maintain data that is added through manual interaction

## Component Definition Generation

Current TODO:
- Ability to detect an existing output file (IE OutputFile flag)
- Ability to retain data in an existing OutputFile on re-generation
  - Easier said than done - check if the control-implementation exists
  - then find the delta of the controls
- Ability to detect an OSCAL manifest file (IE InputFile flag)

## Example 

```bash
./bin/lula generate component -c https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json -r ac-1,ac-3,ac-3.2,ac-4 -o oscal-component.yaml --remarks assessment-objective -l debug
```