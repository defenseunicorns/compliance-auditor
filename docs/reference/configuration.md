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

TODO - description of templating configuration fields

```yaml
# constants = place to define non-changing values that can be of any type
# I think stuff here probably shouldn't be set by env vars - it's hard to be deterministic because of the character set differences, also type differences could lead to weird side effects
# Another note about this - we could probably easily pull in values of child components if this was referenced from a system-level - so this kind of behaves a bit like help values.yaml
constants:
  # map[string]interface{} - elements referenced via template as {{ .const.key }}
  type: software
  title: lula
  # Sample: Istio-specific values
  istio:
    namespace: istio-system # overriden by --set const.istio.namespace=my-istio-namespace
  resources:
    jsoncm: configmaps # (NOT) overriden by LULA_VAR_RESOURCES_JSONCM
    # Problem with this is that json-cm and json_cm are different yaml keys, but would possibly reconcile to the same thing... so you're getting some side effects here that aren't great.
    yamlcm: configmaps
    secret: secrets
    pod: pods
    boolean: false  # (NOT) overriden by LULA_VAR_RESOURCES_BOOLEAN 
    # ok how does this work when they're different types? an env var will always be a string...
  exemptions:
    - one
    - two
    - three

# variables = place to define changing values of string type, and optionally sensitive values
# NOTE - if a variable is defined here, but does not have a default, you will need to make sure it's set either via --set or LULA_VAR_* for the template to execute without error (actually it doesn't error, just prints debug statements)
variables:
  - key: some_lula_secret # set by LULA_VAR_SOME_LULA_SECRET / overriden by --set var.some_lula_secret=my-secret
    default: blahblah  # optional
    sensitive: true # {{ var.some_lula_secret | mask }}
  - key: some_env_var
    default: this-should-be-overridden

# Lula config values, still accessible via LULA_*, where * is the key
log_level: info
target: il5
```