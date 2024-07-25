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