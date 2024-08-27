# Config File

Lula allows the use and specification of a config file in the following ways:
- Checking current working directory for a `lula-config.yaml` file
- Specification with environment variable `LULA_CONFIG=<path>`

## Identification

If identified, Lula will log which configuration file is used to stdout:
```bash
Using config file /home/dev/work/lula/lula-config.yaml
```
## Support

Modification of `log level` can be set in the configuration file by specifying: 

lula-config.yaml
```yaml
log_level: debug
```