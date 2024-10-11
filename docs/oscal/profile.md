# Profile

A [Profile](https://pages.nist.gov/OSCAL/resources/concepts/layer/control/profile/) is an OSCAL model for capturing a baseline of selected controls from one or more catalogs or profiles. In lula, the profile model is available for generation and use as a source to other models that allow for specification of a profile or catalog which represents the source of truth for relevant security controls or other organizational policies. 

## Structure

The primary structure for the Lula production and operations of the `profile` model for determinism is as follows:
- Imports are sorted by `href` in ascending order
- WithIds are sorted by the associated string id in ascending order
- Back Matter Resources are sorted by `title` in ascending order (Case Sensitive Sorting).

### Reproducibility  

The `lula generate` commands are meant to be reproducible. The intent for this generation is to make it easy to update a given model with automation and only inject human intervention as needed. 

For profiles, see the metadata props for the `generation` prop. It should look like the following:
```yaml
props:
  - name: generation
    ns: https://docs.lula.dev/oscal/ns
    value: lula generate profile --source catalog.yaml --include ac-1,ac-2,ac-3
```

>[!NOTE]
>The controls specified for inclusion or exclusion during the generation command are not currently validated to exist in the source artifact. 

## Example

```bash
lula generate profile -s catalog.yaml -i ac-1,ac-2,ac-3
```