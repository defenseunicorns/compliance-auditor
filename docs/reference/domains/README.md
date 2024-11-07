# Domains

Domains, generically, are the areas or categories from which data is collected to be evaluated by a `Lula Validation`. Currently supported Domains are:

* Kubernetes
* API
* File

The domain block of a `Lula Validation` is given as follows, where the sample is indicating a Kubernetes domain is in use:
```yaml
# ... Rest of Lula Validation
domain:
    type: kubernetes   # kubernetes or api accepted
    kubernetes-spec:
        # ... Rest of kubernetes-spec
# ... Rest of Lula Validation
```

Each domain has a particular specification, given by the respective `<domain>-spec` field of the `domain` property of the `Lula Validation`. The sub-pages describe each of these specifications in greater detail.