# Templating

Lula supports composition of both Component Definition and Lula Validation template files. See the [configuration](./configuration.md) documentation for more information on how to configure Lula to use templating. See the [compose CLI command](../cli-commands/lula_tools_compose.md) documentation for more information on the `lula tools compose` command flags to control how templating is applied.

## Component Definition Templating

Component Definition templates can be used to create modular component definitions using values from the `lula-config.yaml` file.

Example:
```yaml
component-definition:
  uuid: E6A291A4-2BC8-43A0-B4B2-FD67CAAE1F8F
  metadata:
    title: {{ .const.title }}
    last-modified: "2022-09-13T12:00:00Z"
    version: "20220913"
    oscal-version: 1.1.2
    parties:
      - uuid: C18F4A9F-A402-415B-8D13-B51739D689FF
        type: organization
        name: Lula Development
        links:
          - href: {{ .const.website }}
            rel: website
```

lula-config.yaml:
```yaml
constants:
  title: Lula Demo
  website: https://github.com/defenseunicorns/lula
```

When this is `composed` with templating applied (`lula tools compose -f <file> --render all`) with the associated `lula-config.yaml`, the resulting component definition will be:

```yaml
component-definition:
  uuid: E6A291A4-2BC8-43A0-B4B2-FD67CAAE1F8F
  metadata:
    title: Lula Demo
    last-modified: "2022-09-13T12:00:00Z"
    version: "20220913"
    oscal-version: 1.1.2
    parties:
      - uuid: C18F4A9F-A402-415B-8D13-B51739D689FF
        type: organization
        name: Lula Development
        links:
          - href: https://github.com/defenseunicorns/lula
            rel: website
```

## Validation Templating

Validation templates can be used to create modular Lula Validations using values from the `lula-config.yaml` file. These can be composed into the component definition using the `lula tools compose` command.

Example:
```yaml
component-definition:
  uuid: E6A291A4-2BC8-43A0-B4B2-FD67CAAE1F8F
  metadata:
    title: Lula Demo
    last-modified: "2022-09-13T12:00:00Z"
    version: "20220913"
    oscal-version: 1.1.2 # This version should remain one version behind latest version for `lula dev upgrade` demo
    parties:
      # Should be consistent across all of the packages, but where is ground truth?
      - uuid: C18F4A9F-A402-415B-8D13-B51739D689FF
        type: organization
        name: Lula Development
        links:
          - href: https://github.com/defenseunicorns/lula
            rel: website
  components:
    - uuid: A9D5204C-7E5B-4C43-BD49-34DF759B9F04
      type: {{ .const.type }}
      title: {{ .const.title }}
      description: |
        Lula - the Compliance Validator
      purpose: Validate compliance controls
      responsible-roles:
        - role-id: provider
          party-uuids:
            - C18F4A9F-A402-415B-8D13-B51739D689FF # matches parties entry for Defense Unicorns
      control-implementations:
        - uuid: A584FEDC-8CEA-4B0C-9F07-85C2C4AE751A
          source: https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json
          description: Validate generic security requirements
          implemented-requirements:
            - uuid: 42C2FFDC-5F05-44DF-A67F-EEC8660AEFFD
              control-id: ID-1
              description: >-
                This control validates that the demo-pod pod in the validation-test namespace contains the required pod label foo=bar in order to establish compliance.
              links:
                - href: "./validation.tmpl.yaml"
                  text: local path template validation
                  rel: lula
```

Where `./validation.tmpl.yaml` is:
```yaml
metadata:
  name: Test validation with templating
  uuid: 99fc662c-109a-4e26-8398-75f3db67f862
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: podvt
        resource-rule:
          name: {{ .const.resources.name }}
          version: v1
          resource: pods
          namespaces: [{{ .const.resources.namespace }}]
provider:
  type: opa
  opa-spec:
    rego: |
      package validate
      import rego.v1

      # Default values
      default validate := false
      default msg := "Not evaluated"

      # Validation result
      validate if {
        { "one", "two", "three" } == { {{ .const.resources.exemptions | concatToRegoList }} }
        "{{ .var.some_env_var }}" == "my-env-var"
        "{{ .var.some_lula_secret }}" == "********"
      }
      msg = validate.msg

      value_of_my_secret := {{ .var.some_lula_secret }}
```

Executing `lula tools compose -f ./component-definition.yaml --render all --render-validations` will result in:

```yaml
component-definition:
  back-matter:
    resources:
      - description: |
          domain:
            kubernetes-spec:
              create-resources: null
              resources:
              - description: ""
                name: podvt
                resource-rule:
                  group: ""
                  name: test-pod-label
                  namespaces:
                  - validation-test
                  resource: pods
                  version: v1
            type: kubernetes
          lula-version: ""
          metadata:
            name: Test validation with templating
            uuid: 99fc662c-109a-4e26-8398-75f3db67f862
          provider:
            opa-spec:
              rego: |
                package validate
                import rego.v1

                # Default values
                default validate := false
                default msg := "Not evaluated"

                # Validation result
                validate if {
                  { "one", "two", "three" } == { "one", "two", "three" }
                  "this-should-be-overridden" == "my-env-var"
                  "" == "********"
                }
                msg = validate.msg

                value_of_my_secret :=
            type: opa
        title: Test validation with templating
        uuid: 99fc662c-109a-4e26-8398-75f3db67f862
  components:
    - control-implementations:
        - description: Validate generic security requirements
          implemented-requirements:
            - control-id: ID-1
              description: This control validates that the demo-pod pod in the validation-test namespace contains the required pod label foo=bar in order to establish compliance.
              links:
                - href: '#99fc662c-109a-4e26-8398-75f3db67f862'
                  rel: lula
                  text: local path template validation
              uuid: 42C2FFDC-5F05-44DF-A67F-EEC8660AEFFD
          source: https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json
          uuid: A584FEDC-8CEA-4B0C-9F07-85C2C4AE751A
      description: |
        Lula - the Compliance Validator
      purpose: Validate compliance controls
      responsible-roles:
        - party-uuids:
            - C18F4A9F-A402-415B-8D13-B51739D689FF
          role-id: provider
      title: lula
      type: software
      uuid: A9D5204C-7E5B-4C43-BD49-34DF759B9F04
  metadata:
    last-modified: XXX
    oscal-version: 1.1.2
    parties:
      - links:
          - href: https://github.com/defenseunicorns/lula
            rel: website
        name: Lula Development
        type: organization
        uuid: C18F4A9F-A402-415B-8D13-B51739D689FF
    title: Lula Demo
    version: "20220913"
  uuid: E6A291A4-2BC8-43A0-B4B2-FD67CAAE1F8F
```

### Composing Validation Templates

If validations are composed into a component definition AND the validation is still intended to be a template, it must be a valid yaml document. For example, the above `validation.tmpl.yaml` is invalid yaml, as the `resource-rule.name` field is not ecapsulated in quotes. A valid yaml version of the above template would be:

```yaml
metadata:
  name: Test validation with templating
  uuid: 99fc662c-109a-4e26-8398-75f3db67f862
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: podvt
        resource-rule:
          name: "{{ .const.resources.name }}"
          version: v1
          resource: pods
          namespaces: ["{{ .const.resources.namespace }}"]
provider:
  type: opa
  opa-spec:
    rego: |
      package validate
      import rego.v1

      # Default values
      default validate := false
      default msg := "Not evaluated"

      # Validation result
      validate if {
        { "one", "two", "three" } == { {{ .const.resources.exemptions | concatToRegoList }} }
        "{{ .var.some_env_var }}" == "my-env-var"
        "{{ .var.some_lula_secret }}" == "********"
      }
      msg = validate.msg

      value_of_my_secret := {{ .var.some_lula_secret }}
```