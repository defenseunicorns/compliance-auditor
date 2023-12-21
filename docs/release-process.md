# Lula Release

This document provides guidance on how to create Lula releases, address release issues, and other helpful tips.

## Creating releases

This project uses [goreleaser](https://github.com/goreleaser/goreleaser-action) for releasing binaries.

### How should I write my commits?

We use conventional commit messages [Conventional Commit messages](https://www.conventionalcommits.org/).

The most important prefixes you should have in mind are:

- `fix:` which represents bug fixes, and correlates to a [SemVer](https://semver.org/)
  patch.
- `feat:` which represents a new feature, and correlates to a SemVer minor.
- `feat!:`,  or `fix!:`, `refactor!:`, etc., which represent a breaking change
  (indicated by the `!`) and will result in a SemVer major.

### How can I influence the version number for a release?

To trigger the goreleaser action you push a signed tag using your GPG key.

```console
git tag -s vX.X.X -m "Release Commit Message"

git push origin vX.X.X
```

Goreleaser can also be triggered by cutting the release manually through GitHub. Following the GitHub [release](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository) documentation and using a semver version number IE. vX.X.X

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

#### Incorrect Tag or Remove a Tag

In the event a release is tagged incorrectly and it needs to be changed this can be completed by someone with admin privileges in the repo by running the following commands:

```console
git tag -d vX.X.X

git tag -s vX.X.X

git push origin vX.X.X --force
```

These commands will delete the tag locally, create a new signed tag version, lastly force push the new tag to origin.

In the event a tag needs deleted and NOT replaced this can be completed by someone with admin privileges in the repo by running the following commands:

```console
git tag -d vX.X.X

git push origin :refs/tags/vX.X.X
```

#### Other issues and helpful tips

- Confirm that the goreleaser configuration is valid by using the [goreleaser cli](https://goreleaser.com/cmd/goreleaser_check/?h=valid)

```sh
goreleaser check .goreleaser.yaml [flags]
```
