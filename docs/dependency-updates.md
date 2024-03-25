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

### Notes
- Validation of the checksums is currently a manual process an a byproduct of not yet capturing the provenance of Renovates checksum process. Given that no single version of Renovate is being used (this is the non-self-hostedGitHub application), we do not track updates to the renovate runtime itself. 