# Lula Release

This document provides guidance on how to create Lula releases, address release issues, and other helpful tips.

## Creating releases

This project uses [goreleaser](https://github.com/goreleaser/goreleaser-action) for releasing binaries and [release-please](https://github.com/marketplace/actions/release-please-action) for creating release PR's.

### How should I write my commits?

We use conventional commit messages [Conventional Commit messages](https://www.conventionalcommits.org/).

The most important prefixes you should have in mind are:

- `fix:` which represents bug fixes, and correlates to a [SemVer](https://semver.org/)
  patch.
- `feat:` which represents a new feature, and correlates to a SemVer minor.
- `feat!:`,  or `fix!:`, `refactor!:`, etc., which represent a breaking change
  (indicated by the `!`) and will result in a SemVer major.

### How can I influence the version number for a release?

PR titles should also follow this pattern and are linted using [commitlint](https://commitlint.js.org/). The PR title will determine the version bump. When a PR is merged (squashed) release-please will kick off a release PR. When that release PR is approved and merged, release-please will create a draft release. Once that draft release is published go-releaser with build and publish the assets, including creating a release in our Homebrew tap repository: [https://github.com/defenseunicorns/homebrew-tap](https://github.com/defenseunicorns/homebrew-tap)

- Pre-v1.0.0 release-please is configured to bump minors on breaking changes and patches otherwise. per [release-please-config](https://github.com/defenseunicorns/lula/blob/main/release-please-config.json)

### How do I fix a release issue?

There are some different ways that something could go wrong with a release, below are some example scenarios and how to fix them.

#### A release is "broken" and should not be used

Rather than removing a release that is deemed to be broken in some fashion, it should be noted as a broken release in the release notes as soon as possible, and a new release created that fixes the issue(s).

The CHANGELOG is not required to be updated, only the release notes must be updated either manually or with CI automation.

- Manual approach: Find the impacted release, edit the release notes, and put this warning at the top:

```md
>[!WARNING]
>PLEASE USE A NEWER VERSION (there are known issues with this release)
```

#### Other issues and helpful tips

- Manual approach: Confirm that the goreleaser configuration is valid by using the [goreleaser cli](https://goreleaser.com/cmd/goreleaser_check/?h=valid).

```sh
goreleaser check .goreleaser.yaml [flags]
```

- Automated approach: On Push and Pull Request the [GoReleaserGitHub Action Workflow](./github/workflows/goreleaser-check.yaml) will run the `goreleaser check` command
