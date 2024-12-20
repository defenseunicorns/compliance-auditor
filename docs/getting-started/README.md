# Getting Started

## Installation

Currently, you can install Lula from Repository Releases, Homebrew, or by building from source. Lula is compatible only with Linux and macOS distributions.

### Repository Release

1. Navigate to the Latest Release Page:
   Open your web browser and go to the following URL to access the latest release of Lula:

   [https://github.com/defenseunicorns/lula/releases/latest](https://github.com/defenseunicorns/lula/releases/latest)

2. Download the Binary:
   On the latest release page, find and download the appropriate binary for your operating system. E.g., `lula_<version>_Linux_amd64`

3. Download the checksums.txt:
   In the list of assets on the release page, locate and download the checksums.txt file. This file contains the checksums for all the binaries in the release.

4. Verify the Download:
   After downloading the binary and checksums.txt, you should verify the integrity of the binary using the checksum provided:
    * Open a terminal and navigate to the directory where you downloaded the binary and checksums.txt.
    * Run the following command to verify the checksum if using Linux:

      ```shell
      sha256sum -c checksums.txt --ignore-missing
      ```

    * Run the following command to verify the checksum if using MacOS:

      ```shell
      shasum -a 256 -c checksums.txt --ignore-missing
      ```

5. On most Linux distributions, install the binary onto your $PATH by moving the downloaded binary to the /usr/local/bin directory:

   ```shell
   sudo mv ./download/path/lula_<version>_Linux_amd64 /usr/local/bin/lula

### Homebrew

1) [Homebrew](https://brew.sh/) is a package manager for macOS and Linux. You can install Lula with Homebrew by running the following:

    ```shell
    brew tap defenseunicorns/tap && brew install lula
    ```

### From Source

1) Clone the repository to your local machine and change into the `lula` directory

    ```shell
    git clone https://github.com/defenseunicorns/lula.git && cd lula
    ```

2) While in the `lula` directory, compile the tool into an executable binary. This outputs the `lula` binary to the `bin` directory.

    ```shell
    make build
    ```

3) On most Linux distributions, install the binary onto your $PATH by moving the downloaded binary to the /usr/local/bin directory:

    ```shell
    sudo mv ./bin/lula /usr/local/bin/lula
    ```

## Quick Start

See the following tutorials for some introductory lessons on how to use Lula. If you are unfamiliar with Lula, the best place to start is the "Simple Demo".

### Tutorials

* [Simple Demo](./simple-demo.md)
* [Using Lula in CI](./lula-in-ci.md)

### Lula Validations

Lula Validation manifests are the underlying mechanisms that dictates the evaluation of a system against a control as resulting in `satisfied` or `not satisfied`. A Lula Validation is linked to a control within a component definition via the OSCAL-specific property, [links](../oscal/oscal-validation-links.md).

Developing Lula Validations can sometimes be more art than science, but generally they should aim to be clear, concise, and robust to system changes. The following guides will help you get started with developing Lula Validations.

* [Develop a Validation](./develop-a-validation.md)
* [Test a Validation](./test-a-validation.md)

See our [references](../reference/README.md) documentation for additional information about Lula Validations.

### Configuration

Lula supports the addition of a configuration file for specifying CLI flags and templating values. See our [configuration](./configuration.md) guide for more information.
