# Providers

Providers are the policy engines which evaluate the input data from the specified domain. Currently supported Providers are:

* [OPA (Open Policy Agent)](opa-provider.md)
* [Kyverno](kyverno-provider.md)

The provider block of a `Lula Validation` is given as follows, where the sample is indicating the OPA provider is in use:
```yaml
# ... Rest of Lula Validation
provider:
    type: opa   # opa or kyverno accepted
    opa-spec:
        # ... Rest of opa-spec
# ... Rest of Lula Validation
```

Each domain specification retreives a specific dataset, and each will return that data to the selected `Provider` in a domain-specific format. However, this data will always take the form of a JSON object when input to a `Provider`. For that reason, it is important that `Domain` and `Provider`specifications are not built wholly independently in a given Validation.
