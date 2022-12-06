# Questions and Answers
This document captures conversations that provide insight into capabiltiies, roles, responsibilities, and functionality of the tool.
These are all subject to change and simply an attempt to document nuances across initial creation as well as provide a central target for high-level discussion.

### Why a compliance engine vs policies in a policy engine?
This solution will provide a mechanism for "inheriting" controls that are produced and versioned controlled in tight coordination with the specific tools that are orchestrated. Policies are highly opinionated in this regard to a specific environment and orchestration, and introducing changes in tools that previously supported control-satisfaction would require refactoring the policy. 

This allows a report (or System Security Plan) to be generated with potential for overlapping control-satisfaction - without being dependent on a single tool.