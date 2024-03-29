# Dependency Updates

## Responsibility
Dependency updates are a responsibility of all project maintainers. All maintainers must be accountable for the updates they introduce and proper review provides a mechanism for reducing the potential for negative impact to the project. 

## Objectives
- Ensuring that all dependencies are updated to their latest versions
- Understanding the implications of updated dependency code
- Validating the provenance of the updated dependency
- Annotation of the reviewed dependency updates

## Guidance

Through the use of [Renovate](https://www.mend.io/renovate/), we can automate the process of updating our dependencies to their latest versions. With this automation comes the responsibility to review considerations and implications to the updates the changes introduce. Review of the dependency updates will begin as renovate creates a pull request with the dependency update. Review should then include the following:

- Review the dependency `release notes` included in the Pull Request
- Compare the source code changes between tagged versions of the dependency
  - Isolate and annotate any potential updates that may impact the project code
  - Review updates for new features or processes that may be positively consumed by the project
- Validate the integrity / provenance of the updated dependency
  - Golang checksums
    - go.mod and project checksums
  - NPM integrity
    - tarball integrity validation
  - Workflow Integrity
    - Tag Commit checksum
- Annotation of the reviewed dependency updates for approval
  - Include any relevant notes or considerations
  - Include steps to validate the dependency updates

## Examples

### Golang Dependencies
Given a dependency in the go.mod file, the following steps are taken to validate the dependency:
1. Identify the dependency source - IE github.com/open-policy-agent/opa
2. Identify the tagged version - IE v0.62.1
3. Curl the checksum of the tagged version - IE curl https://sum.golang.org/lookup/${source}@${tag}

The following should be returned for use in validation:
```
github.com/open-policy-agent/opa v0.62.1 h1:UcxBQ0fe6NEjkYc775j4PWoUFFhx4f6yXKIKSTAuTVk=
github.com/open-policy-agent/opa v0.62.1/go.mod h1:YqiSIIuvKwyomtnnXkJvy0E3KtVKbavjPJ/hNMuOmeM=
```

### NPM Dependencies
Given a dependency in the package-lock.json file, the following steps are taken to validate the dependency:
1. Identify the resolved dependency archive - IE https://registry.npmjs.org/autoprefixer/-/autoprefixer-10.4.19.tgz
2. Curl the archive locally - IE curl -LO https://registry.npmjs.org/autoprefixer/-/autoprefixer-10.4.19.tgz
3. Perform a SHA512 hash on the archive - IE cat ./autoprefixer-10.4.19.tgz | openssl dgst -sha512 -binary | openssl base64 -A
4. Compare results against the `integrity` field in the package-lock.json file

### GitHub Action Dependencies
Review the updated commit hash of the tagged action against the tag of the action in the source git repository. 

### Notes
- Validation of the checksums is currently a manual process and a byproduct of not yet capturing the provenance of Renovates checksum process. Given that no single version of Renovate is being used (this is the non-self-hosted GitHub application), we do not track updates to the renovate runtime itself. 