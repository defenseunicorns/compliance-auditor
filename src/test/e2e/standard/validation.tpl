domain:
  type: api
  api-spec:
    requests:
    - name: local
      url: https://some.url/v1/api
      content-type: application/json
      headers:
        Authorization: Bearer {{ .var.some_lula_secret }}
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
        input.jsoncm.name == {{ .const.resources.jsoncm }}
        input.yamlcm.logging.level == {{ .const.resources.yamlcm }}
        not input.secret.name in { {{ .const.exemptions | concatToRegoList }} } # tmpl creates: "one", "two", "three"
      }
      msg = validate.msg

      test_env_var := {{ .var.some_env_var }} # non-sensitive
      test_another_env_var := {{ .var.another_env_var }} # non-sensitive, no default