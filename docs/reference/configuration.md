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