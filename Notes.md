## Standards Optionality

### targeted selection

- Create a flag that allows specification of a control-implementation OR framework prop
- Identify / collect each specified target within 1 -> N components
- Perform validation

### default

- Iterate through all components
- map validations to implemented-requirements if exists
- For each control-implementation/framework
  - validate and produce findings/observations
  - Convert those findings/observations into a result
  - append result to a list of all results
- Produce assessment-results
  - merge if required

### Current Workflow
- Consume component-definition
- Compose the component definition
- Create validation store
- Create the requirements store
   - This is done 1-> n components / 1-> control-implementations / 1 -> n implemented-requirements


### Common Workflows

- Consume a component definition
- Compose the component definition
- Create the requirements store
  - This is done 1-> n components / 1-> control-implementations / 1 -> n implemented-requirements

### Ideas

- Add an additional layer where we potentially identify which control-implementations are required
- new function
  - Build requirements store based on 1->N control-implementation objects
  - return and build a result object