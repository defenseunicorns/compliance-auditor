# Generation

The `generation` prop is an identifier for the purposes of tracking imperative reproducibility of a given artifact or subset of an artifact. In the example below, the `lula generate component` command annotates how a given `control-implementation` - and associated `component` were generated. 

```yaml
props:
  - name: generation
    ns: https://docs.lula.dev/oscal/ns
    value: lula generate component --catalog-source https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json --component 'Component Title' --requirements ac-1,ac-3,ac-3.2,ac-4 --remarks assessment-objective
```
