# Validation

Logic to extract, execute, and evaluate Lula validations on a specified environment against targeted requirements.

## Components

### Producer

The producer is responsible for providing the structure which stores Lula validations and associated requirements. Current producers are:
- OSCAL Componet Definition
- Simple (Lula Validation .yaml file)

### Consumer

The consumer is responsible for evaluating the validations and requirements from the producer and providing the results in the custom format. Current consumers are:
- OSCAL Assessment Results
- Simple (direct pass/fail from requirements)

### Requirement

The requirement is responsible for providing the structure which stores producer-generated requirements along with any associated validations. Current requirements are:
- Component Definition Requirements
- Simple (simple requirement to track aggregated validation pass/fail)

