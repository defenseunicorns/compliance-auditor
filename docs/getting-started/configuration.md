# Configuration

Lula allows the use and specification of a config file in the following ways:
- Checking current working directory for a `lula-config.yaml` file
- Specification with environment variable `LULA_CONFIG=<path>`

Environment Variables can be used to specify configuration values through use of `LULA_<VAR>` -> Example: `LULA_TARGET=il5` 

## Identification

If identified, Lula will log which configuration file is used to stdout:
```bash
Using config file /home/dev/work/lula/lula-config.yaml
```

## Precedence

The precedence for configuring settings, such as `target`, follows this hierarchy:

### **Command Line Flag > Environment Variable > Configuration File**

1. **Command Line Flag:**  
   When a setting like `target` is specified using a command line flag, this value takes the highest precedence, overriding any environment variable or configuration file settings.

2. **Environment Variable:**  
   If the setting is not provided via a command line flag, an environment variable (e.g., `export LULA_TARGET=il5`) will take precedence over the configuration file.

3. **Configuration File:**  
   In the absence of both a command line flag and environment variable, the value specified in the configuration file will be used. This will override system defaults.

## Support

Modification of command variables can be set in the configuration file:

lula-config.yaml
```yaml
log_level: debug
target: il4
summary: true
```

### Templating Configuration Fields

Templating values are set in the configuration file via the use of `constants` and `variables` fields.

#### Constants

A sample `constants` section of a `lula-config.yaml` file is as follows:

```yaml
constants:
  type: software
  title: lula

  resources:
    name: test-pod-label
    namespace: validation-test
    imagelist:
      - nginx
      - nginx2
```

Constants will respect the structure of a map[string]interface{} and can be referenced as follows:

```yaml
# validaiton.yaml
metadata:
  name: sample {{ .const.type }} validation for {{ .const.title }}
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: myPod
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

      validate if {
        input.myPod.metadata.name == {{ .const.resources.name }}
        input.myPod.containers[_].name in { {{ .const.resources.imagelist | concatToRegoList }} }
      }
```

And will be rendered as:
```yaml
metadata:
  name: sample software validation for lula
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: myPod
        resource-rule:
          name: myPod
          version: v1
          resource: pods
          namespaces: [validation-test]
provider:
  type: opa
  opa-spec:
    rego: |
      package validate
      import rego.v1

      validate if {
        input.myPod.metadata.name == "myPod"
        input.myPod.containers[_].image in { "nginx", "nginx2" }
      }
```

The constant's keys should be in the format `.const.<key>` and should not contain any '-' or '.' characters, as this will not respect the go text/template format. 

> [!IMPORTANT]
> Due to viper limitations, all constants should be referenced in the template as lowercase values.

#### Variables

A sample `variables` section of a `lula-config.yaml` file is as follows:

```yaml
variables:
  - key: some_lula_secret
    sensitive: true
  - key: some_env_var
    default: this-should-be-overridden
```

The `variables` section is a list of `key`, `default`, and `sensitive` fields, where `sensitive` and `default` are optional. The `key` and `default` fields are strings, and the `sensitive` field is a boolean.

A default value can be specified in the case where an environment variable may or may not be set, however an environment variable will always take precedence over a default value.

The environment variable should follow the pattern of `LULA_VAR_<key>` (not case sensitive), where `<key>` is the key specified in the `variables` section.

When using `sensitive` variables, the default behavior is to mask the value in the output of the template.