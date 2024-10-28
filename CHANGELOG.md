# Changelog

## [0.10.0](https://github.com/defenseunicorns/lula/compare/v0.9.1...v0.10.0) (2024-10-25)


### Features

* **console:** add support for multiple input files ([#729](https://github.com/defenseunicorns/lula/issues/729)) ([103ca0d](https://github.com/defenseunicorns/lula/commit/103ca0deabbb5db7fc90512c82817cb91bdd4c3f))
* **console:** assessment results generation ([#744](https://github.com/defenseunicorns/lula/issues/744)) ([6cb5933](https://github.com/defenseunicorns/lula/commit/6cb5933828f137d2ee024427f76885d12f5ad4bc))
* **console:** lula validation view ([#727](https://github.com/defenseunicorns/lula/issues/727)) ([481648f](https://github.com/defenseunicorns/lula/commit/481648fb366da9abf19441d39b81173935424704))
* **generate:** support for profile model and basic generation ([#694](https://github.com/defenseunicorns/lula/issues/694)) ([cb4fc6f](https://github.com/defenseunicorns/lula/commit/cb4fc6f2282547352c0a3a88ddd135c5a86e58eb))


### Miscellaneous

* address linter report - swallowed errors and minor test cleanup of nil checks ([#740](https://github.com/defenseunicorns/lula/issues/740)) ([05a7f6e](https://github.com/defenseunicorns/lula/commit/05a7f6e075adafeb2b3c635803cc12083f12c01d))
* **build:** add golangci-lint to the pipeline as a non-blocking step ([#742](https://github.com/defenseunicorns/lula/issues/742)) ([fedb0c9](https://github.com/defenseunicorns/lula/commit/fedb0c9731656941b24c9d2cf6fb85bdc4604002))
* **deps:** update actions/checkout action to v4.2.2 ([#756](https://github.com/defenseunicorns/lula/issues/756)) ([7e1a193](https://github.com/defenseunicorns/lula/commit/7e1a1937bc69b99de67c50185b5c3c3a87cb4395))
* **deps:** update actions/setup-go action to v5.1.0 ([#761](https://github.com/defenseunicorns/lula/issues/761)) ([e464cd0](https://github.com/defenseunicorns/lula/commit/e464cd04092ce78b50c750933edcc85c7f4fecdc))
* **deps:** update actions/setup-node action to v4.1.0 ([#759](https://github.com/defenseunicorns/lula/issues/759)) ([2fd7f86](https://github.com/defenseunicorns/lula/commit/2fd7f86fc4f4723e34a3818fdd4e6282ee4cfed5))
* **deps:** update anchore/sbom-action action to v0.17.4 ([#731](https://github.com/defenseunicorns/lula/issues/731)) ([36ef42a](https://github.com/defenseunicorns/lula/commit/36ef42aefb9fa5d48997a7f4b7d342c6a0d92cd3))
* **deps:** update anchore/sbom-action action to v0.17.5 ([#753](https://github.com/defenseunicorns/lula/issues/753)) ([521452c](https://github.com/defenseunicorns/lula/commit/521452cf57370bcc4521020736b52438a5df797d))
* **deps:** update github/codeql-action action to v3.26.13 ([#737](https://github.com/defenseunicorns/lula/issues/737)) ([86c9376](https://github.com/defenseunicorns/lula/commit/86c93762fcfc9eb819989dd0e376a798a6add4b2))
* **deps:** update github/codeql-action action to v3.27.0 ([#745](https://github.com/defenseunicorns/lula/issues/745)) ([c87621e](https://github.com/defenseunicorns/lula/commit/c87621e40c2203cac8b4273d2291beab8765d8b2))
* **deps:** update kubernetes packages to v0.31.2 ([#757](https://github.com/defenseunicorns/lula/issues/757)) ([259b180](https://github.com/defenseunicorns/lula/commit/259b18071651db2efd574cfa791d14b1f81628f3))
* **deps:** update module github.com/charmbracelet/bubbletea to v1.1.2 ([#763](https://github.com/defenseunicorns/lula/issues/763)) ([83c9c3a](https://github.com/defenseunicorns/lula/commit/83c9c3aa841534ff785c138a9247b867995dbc98))
* **deps:** update module github.com/charmbracelet/lipgloss to v0.13.1 ([#755](https://github.com/defenseunicorns/lula/issues/755)) ([821643d](https://github.com/defenseunicorns/lula/commit/821643d005b8c5480ef6bbb36d8961b8e6d4a352))
* **deps:** update module github.com/open-policy-agent/conftest to v0.56.0 ([#743](https://github.com/defenseunicorns/lula/issues/743)) ([c8b2293](https://github.com/defenseunicorns/lula/commit/c8b2293eaa3594acf6b867252020a2a513d820e3))
* **deps:** update module sigs.k8s.io/e2e-framework to v0.5.0 ([#754](https://github.com/defenseunicorns/lula/issues/754)) ([1f7f3c7](https://github.com/defenseunicorns/lula/commit/1f7f3c72d06acd6cb394e565abd69760dffba93d))
* **lint:** appease the linter ([#746](https://github.com/defenseunicorns/lula/issues/746)) ([bcb2ab0](https://github.com/defenseunicorns/lula/commit/bcb2ab0b7d150e2d908fe2a9d674663f3c040489))
* **scan:** integrate KICS scanning into PR pipelines ([#751](https://github.com/defenseunicorns/lula/issues/751)) ([0112462](https://github.com/defenseunicorns/lula/commit/0112462ec52cec8150b559f2745f2974bbd18de6))
* **scanning:** create gosec integration in pipelines ([#739](https://github.com/defenseunicorns/lula/issues/739)) ([41bce03](https://github.com/defenseunicorns/lula/commit/41bce03649ce3cbb578c0bcb2d223884348c86d5))

## [0.9.1](https://github.com/defenseunicorns/lula/compare/v0.9.0...v0.9.1) (2024-10-12)


### Bug Fixes

* **release:** add environment to push job ([#735](https://github.com/defenseunicorns/lula/issues/735)) ([1ed52f1](https://github.com/defenseunicorns/lula/commit/1ed52f1cf36214ef20ac5a95cee5d5f266232192))

## [0.9.0](https://github.com/defenseunicorns/lula/compare/v0.8.0...v0.9.0) (2024-10-11)


### ⚠ BREAKING CHANGES

* **kubernetes:** wait logic kubernetes version support ([#718](https://github.com/defenseunicorns/lula/issues/718))

### Features

* **compose:** template files during compose operations ([#686](https://github.com/defenseunicorns/lula/issues/686)) ([c1745a4](https://github.com/defenseunicorns/lula/commit/c1745a41ff15b9cf8d6f5c4bf459be88bc84cbf9))
* **domains:** file domain ([#703](https://github.com/defenseunicorns/lula/issues/703)) ([bd4f577](https://github.com/defenseunicorns/lula/commit/bd4f57778c5e5bac539d14955e594ee15312c39c))
* **file domain:** add support for reading arbitrary files as strings ([#726](https://github.com/defenseunicorns/lula/issues/726)) ([0b1c0c8](https://github.com/defenseunicorns/lula/commit/0b1c0c8ddf7c0f5de8e23a0b42ca2348efaaef78))
* **kubernetes:** support running both create resources and resources in the kubernetes spec ([#714](https://github.com/defenseunicorns/lula/issues/714)) ([6839d20](https://github.com/defenseunicorns/lula/commit/6839d205ea0f4434d6af2071f3f3ed444b131944))
* **kubernetes:** wait logic kubernetes version support ([#718](https://github.com/defenseunicorns/lula/issues/718)) ([cc06251](https://github.com/defenseunicorns/lula/commit/cc06251e75facf6f321ad4ca2f8609f782dcfb29))
* **release:** add brew install for lula ([#707](https://github.com/defenseunicorns/lula/issues/707)) ([fd1d3e0](https://github.com/defenseunicorns/lula/commit/fd1d3e08754a845e25c849b280ed6390a377e138))
* **validate:** template oscal during runtime ([#708](https://github.com/defenseunicorns/lula/issues/708)) ([3f5a110](https://github.com/defenseunicorns/lula/commit/3f5a110ecf692d99e1511ac82b737d82764321c2))


### Bug Fixes

* add goreleaser pin version annotate ([#712](https://github.com/defenseunicorns/lula/issues/712)) ([68bc101](https://github.com/defenseunicorns/lula/commit/68bc1014edb701da12ddde6ae83ba90c8e19e774))
* **composition:** nil pointer in composition ([#733](https://github.com/defenseunicorns/lula/issues/733)) ([8ad4209](https://github.com/defenseunicorns/lula/commit/8ad420970cd6bd72ee0c18e6c25a4578e9db4432))
* **console:** refactor, retries, sleep to address flaky tests ([#698](https://github.com/defenseunicorns/lula/issues/698)) ([02101a5](https://github.com/defenseunicorns/lula/commit/02101a5633c009ff46083651745b6aa40ac62448))
* **console:** reset compdef when editing ([#701](https://github.com/defenseunicorns/lula/issues/701)) ([4e25f01](https://github.com/defenseunicorns/lula/commit/4e25f014d8ba9bd88df3317ec51ce3fa783203d0))
* **read:** error checking prior to file writes ([#687](https://github.com/defenseunicorns/lula/issues/687)) ([1ab0eef](https://github.com/defenseunicorns/lula/commit/1ab0eefdeeb1d59f16f33249b1a6fce141ef5942))


### Miscellaneous

* add global command context for program cancelation and everything else ([#696](https://github.com/defenseunicorns/lula/issues/696)) ([df81cf7](https://github.com/defenseunicorns/lula/commit/df81cf7a74e6f78c27055b82c20375f53976cea8))
* **deps:** update actions/checkout action to v4.2.1 ([#713](https://github.com/defenseunicorns/lula/issues/713)) ([802601a](https://github.com/defenseunicorns/lula/commit/802601a70fadfc142a47cc6e8528478a6aac3291))
* **deps:** update actions/upload-artifact action to v4.4.3 ([#711](https://github.com/defenseunicorns/lula/issues/711)) ([a954664](https://github.com/defenseunicorns/lula/commit/a954664d0b2e25d58097425dfbeac193a200b6c5))
* **deps:** update github/codeql-action action to v3.26.12 ([#691](https://github.com/defenseunicorns/lula/issues/691)) ([0efb120](https://github.com/defenseunicorns/lula/commit/0efb120a6f50e650a5e2962125a7495a21236fb8))
* **deps:** update module github.com/open-policy-agent/opa to v0.69.0 ([#692](https://github.com/defenseunicorns/lula/issues/692)) ([e08d695](https://github.com/defenseunicorns/lula/commit/e08d695ea6629e2c60a33ae85edf076bbb49ee68))
* **deps:** update module sigs.k8s.io/cli-utils to v0.37.2 ([#721](https://github.com/defenseunicorns/lula/issues/721)) ([5fd0f32](https://github.com/defenseunicorns/lula/commit/5fd0f3244e5543e5302fce2ea4a42afc87026217))
* update getting started doc to include brew install ([#720](https://github.com/defenseunicorns/lula/issues/720)) ([26c3f8d](https://github.com/defenseunicorns/lula/commit/26c3f8dd1d9a5e31d7bf3936b453a3e0edfd2755))

## [0.8.0](https://github.com/defenseunicorns/lula/compare/v0.7.0...v0.8.0) (2024-09-27)


### ⚠ BREAKING CHANGES

* **template:** introducing variables and sensitive configuration ([#672](https://github.com/defenseunicorns/lula/issues/672))

### Features

* **console:** editing a component definition ([#648](https://github.com/defenseunicorns/lula/issues/648)) ([ae06e27](https://github.com/defenseunicorns/lula/commit/ae06e27869043270647670693df342710e3d4390))
* **template:** enable remote file templating ([#680](https://github.com/defenseunicorns/lula/issues/680)) ([f16bcf6](https://github.com/defenseunicorns/lula/commit/f16bcf64134ab3eda904b40d26e72c19cd96be9b))
* **template:** introducing variables and sensitive configuration ([#672](https://github.com/defenseunicorns/lula/issues/672)) ([5d1f232](https://github.com/defenseunicorns/lula/commit/5d1f23257ba7f11508a90c883b152349bcc2d7fd))
* **validate:** save validation resources ([#612](https://github.com/defenseunicorns/lula/issues/612)) ([7b9a771](https://github.com/defenseunicorns/lula/commit/7b9a771852349903025d5d733f0d71fab5133daa))


### Bug Fixes

* cleaned whitespace+newline in rego ([#671](https://github.com/defenseunicorns/lula/issues/671)) ([ac7039d](https://github.com/defenseunicorns/lula/commit/ac7039d2222177869e4cf4db544b90a762aa1a0c))
* trim whitespace bug ([#677](https://github.com/defenseunicorns/lula/issues/677)) ([e30a824](https://github.com/defenseunicorns/lula/commit/e30a8247123ea4bbdf0a582964dfe4ff81aac9f1))


### Miscellaneous

* **codeowners:** update codeowners to reflect current team ([#663](https://github.com/defenseunicorns/lula/issues/663)) ([7fceaf6](https://github.com/defenseunicorns/lula/commit/7fceaf67145c38933e2f8b61177e31ff7c8a84e2))
* **deps:** update actions/checkout action to v4.2.0 ([#681](https://github.com/defenseunicorns/lula/issues/681)) ([187b8a2](https://github.com/defenseunicorns/lula/commit/187b8a2da0545fc78ba56f051bfd6bd19583f3ce))
* **deps:** update actions/github-script digest to 660ec11 ([#669](https://github.com/defenseunicorns/lula/issues/669)) ([ea40e70](https://github.com/defenseunicorns/lula/commit/ea40e70cd84d3cfd1889c9b0d2e27b49d171ce44))
* **deps:** update actions/setup-node action to v4.0.4 ([#674](https://github.com/defenseunicorns/lula/issues/674)) ([643d502](https://github.com/defenseunicorns/lula/commit/643d502278a187a90c643bb76c50373f2c7d6117))
* **deps:** update github.com/charmbracelet/x/exp/teatest digest to 227168d ([#666](https://github.com/defenseunicorns/lula/issues/666)) ([6bc23e3](https://github.com/defenseunicorns/lula/commit/6bc23e3109d6e415668209ca3dfc59064fd019f1))
* **deps:** update github/codeql-action action to v3.26.8 ([#673](https://github.com/defenseunicorns/lula/issues/673)) ([0ca43a1](https://github.com/defenseunicorns/lula/commit/0ca43a1570867b2d8d49429d92cc18b30bbfc26c))
* **deps:** update github/codeql-action action to v3.26.9 ([#679](https://github.com/defenseunicorns/lula/issues/679)) ([20bdbcd](https://github.com/defenseunicorns/lula/commit/20bdbcd80ad877bac149d249c4e931eb1fc43e33))

## [0.7.0](https://github.com/defenseunicorns/lula/compare/v0.6.0...v0.7.0) (2024-09-13)


### Features

* **config:** support for target/summary in lula config ([#640](https://github.com/defenseunicorns/lula/issues/640)) ([28ce6e5](https://github.com/defenseunicorns/lula/commit/28ce6e57b8ce9490e6c39f9c0f6c6c8df77a5a74))
* **template:** template command with initial docs ([#644](https://github.com/defenseunicorns/lula/issues/644)) ([89be460](https://github.com/defenseunicorns/lula/commit/89be4609caa7fcd1e39036d12a67b0c7a72ea97f))
* yaml map injection ([#568](https://github.com/defenseunicorns/lula/issues/568)) ([3babbc8](https://github.com/defenseunicorns/lula/commit/3babbc8c41047c27b5586744e4f3f5f1c19ff1ce))


### Bug Fixes

* **console:** console testing ([#629](https://github.com/defenseunicorns/lula/issues/629)) ([78e4ae9](https://github.com/defenseunicorns/lula/commit/78e4ae9f930c0ef9b47f65bea1c2a54906717ae3))
* **gen-cli-docs:** fix path for gen-cli-docs ([#646](https://github.com/defenseunicorns/lula/issues/646)) ([6ec6e6f](https://github.com/defenseunicorns/lula/commit/6ec6e6fb39d7d22366847b781891d0141f069331))
* **test:** decouple unit from e2e test workflows ([#662](https://github.com/defenseunicorns/lula/issues/662)) ([a4097a1](https://github.com/defenseunicorns/lula/commit/a4097a1fbc2e01a0feb37f368fc126b0de7f5e2e))


### Miscellaneous

* **deps:** update actions/github-script digest to 58d7008 ([#660](https://github.com/defenseunicorns/lula/issues/660)) ([642cfa3](https://github.com/defenseunicorns/lula/commit/642cfa3474848297481a58969717c41cbf53dc42))
* **deps:** update actions/upload-artifact action to v4.4.0 ([#635](https://github.com/defenseunicorns/lula/issues/635)) ([b342f63](https://github.com/defenseunicorns/lula/commit/b342f63c06be37ff2dda332bc83bc561b22f054f))
* **deps:** update commitlint monorepo to v19.5.0 ([#656](https://github.com/defenseunicorns/lula/issues/656)) ([031b524](https://github.com/defenseunicorns/lula/commit/031b5245848c393f4e6d40cbd749e6b3e3fd01f3))
* **deps:** update github.com/charmbracelet/x/exp/teatest digest to 162f303 ([#643](https://github.com/defenseunicorns/lula/issues/643)) ([4323cb0](https://github.com/defenseunicorns/lula/commit/4323cb08671ee36f16bb0f0dfe876a9dc6dc6186))
* **deps:** update github.com/charmbracelet/x/exp/teatest digest to 9ef7ff4 ([#655](https://github.com/defenseunicorns/lula/issues/655)) ([27bff1c](https://github.com/defenseunicorns/lula/commit/27bff1c4c7f114675ff7375b59f449fc43c5b9f4))
* **deps:** update github/codeql-action action to v3.26.7 ([#659](https://github.com/defenseunicorns/lula/issues/659)) ([51025b4](https://github.com/defenseunicorns/lula/commit/51025b427123964b0712283634c58debba0543ab))
* **deps:** update kubernetes packages to v0.31.1 ([#658](https://github.com/defenseunicorns/lula/issues/658)) ([fbe7b8f](https://github.com/defenseunicorns/lula/commit/fbe7b8f616a682f1099476b0d929a24a267e4eeb))
* **deps:** update module github.com/charmbracelet/bubbles to v0.20.0 ([#630](https://github.com/defenseunicorns/lula/issues/630)) ([4f1f3de](https://github.com/defenseunicorns/lula/commit/4f1f3debfd050269776767ae69290d51b2fc2e0e))
* **deps:** update module github.com/charmbracelet/bubbletea to v1.1.1 ([#633](https://github.com/defenseunicorns/lula/issues/633)) ([394b48d](https://github.com/defenseunicorns/lula/commit/394b48d19a6d56c9ac196e7717ec37e55ad7f606))
* **deps:** update module sigs.k8s.io/kustomize/kyaml to v0.17.2 ([#638](https://github.com/defenseunicorns/lula/issues/638)) ([c66df44](https://github.com/defenseunicorns/lula/commit/c66df449358a00f9c7dbf0b0c71dacb56f229e44))

## [0.6.0](https://github.com/defenseunicorns/lula/compare/v0.5.1...v0.6.0) (2024-08-31)


### Features

* **configuration:** add initial support for Viper command initialization ([#607](https://github.com/defenseunicorns/lula/issues/607)) ([2c94c83](https://github.com/defenseunicorns/lula/commit/2c94c8312233a40a536d0f98ca513e2bbe8bf720))
* **console:** initial tui for component-definition read ([#608](https://github.com/defenseunicorns/lula/issues/608)) ([a0338af](https://github.com/defenseunicorns/lula/commit/a0338affd4934c5100ffc280b6fac016b9eb0b91))


### Bug Fixes

* **oscal:** ensure component definition UUID is updated on modification ([#615](https://github.com/defenseunicorns/lula/issues/615)) ([5516482](https://github.com/defenseunicorns/lula/commit/55164824f9f16c0c76c3a068d68c42aae80bd0f3))


### Miscellaneous

* **deps:** update anchore/sbom-action action to v0.17.2 ([#613](https://github.com/defenseunicorns/lula/issues/613)) ([4fb9090](https://github.com/defenseunicorns/lula/commit/4fb909046e769cacf528c88fd46d957c804affd2))
* **deps:** update commitlint monorepo to v19.4.1 ([#619](https://github.com/defenseunicorns/lula/issues/619)) ([0adb2d0](https://github.com/defenseunicorns/lula/commit/0adb2d0093b035f60af8794722a354128bbc2446))
* **deps:** update github/codeql-action action to v3.26.2 ([#600](https://github.com/defenseunicorns/lula/issues/600)) ([6601566](https://github.com/defenseunicorns/lula/commit/66015662002d17a7c3182f0c5d4565a50ece8093))
* **deps:** update github/codeql-action action to v3.26.4 ([#611](https://github.com/defenseunicorns/lula/issues/611)) ([dd592ce](https://github.com/defenseunicorns/lula/commit/dd592ceabc834eeb0a453bc8a0a78b04b4576bfb))
* **deps:** update github/codeql-action action to v3.26.5 ([#616](https://github.com/defenseunicorns/lula/issues/616)) ([aa4e122](https://github.com/defenseunicorns/lula/commit/aa4e1220c725002e448707f00c456dd09161254e))
* **deps:** update github/codeql-action action to v3.26.6 ([#626](https://github.com/defenseunicorns/lula/issues/626)) ([488cd91](https://github.com/defenseunicorns/lula/commit/488cd915f7801dddd427edf56bb8effd2d3bbc92))
* **deps:** update kubernetes packages to v0.31.0 ([#594](https://github.com/defenseunicorns/lula/issues/594)) ([6dd7463](https://github.com/defenseunicorns/lula/commit/6dd7463c724c91f3ec34ddca87605481bdd8356f))
* **deps:** update module github.com/open-policy-agent/opa to v0.68.0 ([#628](https://github.com/defenseunicorns/lula/issues/628)) ([914b1f2](https://github.com/defenseunicorns/lula/commit/914b1f2717dd30d6c8b2c878da0ec1582e7b047c))
* **deps:** update module github.com/spf13/viper to v1.19.0 ([#618](https://github.com/defenseunicorns/lula/issues/618)) ([22fd668](https://github.com/defenseunicorns/lula/commit/22fd668d81923f1f098c919d0472637aa63bb451))
* **docs:** add docs generation command ([#606](https://github.com/defenseunicorns/lula/issues/606)) ([4dd0450](https://github.com/defenseunicorns/lula/commit/4dd0450fa5f91cf4d1b3dceee52bde8f47f50935))
* **docs:** adding 2024 roadmap to the project documentation ([#320](https://github.com/defenseunicorns/lula/issues/320)) ([120d15a](https://github.com/defenseunicorns/lula/commit/120d15a77dbcdc553c117245d392d4dcc3238066))
* **docs:** update CODEOWNERS ([#636](https://github.com/defenseunicorns/lula/issues/636)) ([77e60fe](https://github.com/defenseunicorns/lula/commit/77e60fe6ecf13e71acc035139904e1c9aa7ab3b4))

## [0.5.1](https://github.com/defenseunicorns/lula/compare/v0.5.0...v0.5.1) (2024-08-16)


### Bug Fixes

* **docs:** updated namespace doc locations and associated prop namespace url ([#602](https://github.com/defenseunicorns/lula/issues/602)) ([ee7df0b](https://github.com/defenseunicorns/lula/commit/ee7df0bb1b9f8dfb92d788fd026d7f2002fcc664))


### Miscellaneous

* **deps:** update anchore/sbom-action action to v0.17.1 ([#593](https://github.com/defenseunicorns/lula/issues/593)) ([ef8b546](https://github.com/defenseunicorns/lula/commit/ef8b5465fd57dcbc18334cd366ad3ae21f338ef7))
* **deps:** update github/codeql-action action to v3.26.1 ([#595](https://github.com/defenseunicorns/lula/issues/595)) ([a34281c](https://github.com/defenseunicorns/lula/commit/a34281ca1bc22561e65edca1b3f94892140fb56e))
* **docs:** update issue template with expected deliverables line item ([#590](https://github.com/defenseunicorns/lula/issues/590)) ([97af4ba](https://github.com/defenseunicorns/lula/commit/97af4ba015529e71e501f624affdf6a232707b68))
* removed renovate schedule ([#603](https://github.com/defenseunicorns/lula/issues/603)) ([55bbdf7](https://github.com/defenseunicorns/lula/commit/55bbdf77941a29229cf91b7cfa65a592410f2c23))

## [0.5.0](https://github.com/defenseunicorns/lula/compare/v0.4.5...v0.5.0) (2024-08-09)


### ⚠ BREAKING CHANGES

* **validation-result:** bump to go-oscal v0.6.0 jsonschema v6 update ([#544](https://github.com/defenseunicorns/lula/issues/544))

### Features

* **validation-result:** bump to go-oscal v0.6.0 jsonschema v6 update ([#544](https://github.com/defenseunicorns/lula/issues/544)) ([5e75714](https://github.com/defenseunicorns/lula/commit/5e75714172f5b72d9c7c346011d086f5fc2c790c))


### Bug Fixes

* **generate:** component generation from catalog error handling ([#573](https://github.com/defenseunicorns/lula/issues/573)) ([371d54d](https://github.com/defenseunicorns/lula/commit/371d54d91e9f8dc4561cf161ad7a7ff4a7efb5f3))
* **generate:** proposed transition of generation annotation to props ([#574](https://github.com/defenseunicorns/lula/issues/574)) ([b7a936d](https://github.com/defenseunicorns/lula/commit/b7a936df536fc8ccf9c22af8bafcd1e4e05e19d9))
* **validate:** get non-namespace scoped resources ([#585](https://github.com/defenseunicorns/lula/issues/585)) ([a5b8857](https://github.com/defenseunicorns/lula/commit/a5b8857508d1271cac46c6587f43c17075d8b590))


### Miscellaneous

* **deps:** update actions/upload-artifact action to v4.3.6 ([#575](https://github.com/defenseunicorns/lula/issues/575)) ([27b2e8a](https://github.com/defenseunicorns/lula/commit/27b2e8a41827dea765db134e4cb9e462e8f3c19a))
* **deps:** update dependency commitlint to v19.4.0 ([#583](https://github.com/defenseunicorns/lula/issues/583)) ([732b22b](https://github.com/defenseunicorns/lula/commit/732b22b3d4381e6f6dd2bd615e5c395b34458233))
* **deps:** update github/codeql-action action to v3.26.0 ([#582](https://github.com/defenseunicorns/lula/issues/582)) ([8a92a8c](https://github.com/defenseunicorns/lula/commit/8a92a8c225d31e0e14e0d1578c5bb0acff18014d))
* **deps:** update module github.com/open-policy-agent/opa to v0.67.1 ([#577](https://github.com/defenseunicorns/lula/issues/577)) ([563e893](https://github.com/defenseunicorns/lula/commit/563e893fb82eac803bc19404d65dba278b508760))
* **docs:** initial SSP generation research docs ([#548](https://github.com/defenseunicorns/lula/issues/548)) ([0891508](https://github.com/defenseunicorns/lula/commit/0891508453a4cc79dd0462f339cb0e07b840b628))

## [0.4.5](https://github.com/defenseunicorns/lula/compare/v0.4.4...v0.4.5) (2024-08-02)


### Bug Fixes

* **release:** configuration to bump minor version when including features ([#576](https://github.com/defenseunicorns/lula/issues/576)) ([6bd11bb](https://github.com/defenseunicorns/lula/commit/6bd11bb55b1159c5ad73cb9314b13a0b51a08efe))
* **validate:** allow for optionality among potential standards in a component definition ([#532](https://github.com/defenseunicorns/lula/issues/532)) ([ac0befb](https://github.com/defenseunicorns/lula/commit/ac0befb872f0b634778bf0a6c1f731620dd2e1a2))


### Miscellaneous

* **deps:** update actions/github-script digest to 35b1cdd ([#570](https://github.com/defenseunicorns/lula/issues/570)) ([4a2b03b](https://github.com/defenseunicorns/lula/commit/4a2b03b9eaf1acdf731f83480617663346f82ed6))
* **deps:** update ossf/scorecard-action action to v2.4.0 ([#565](https://github.com/defenseunicorns/lula/issues/565)) ([2207d71](https://github.com/defenseunicorns/lula/commit/2207d71805401dda8721a93a2e1d578d1771801c))

## [0.4.4](https://github.com/defenseunicorns/lula/compare/v0.4.3...v0.4.4) (2024-07-26)


### Features

* **evaluate:** add observation summary ([#540](https://github.com/defenseunicorns/lula/issues/540)) ([8a07833](https://github.com/defenseunicorns/lula/commit/8a07833c5a563d8e857515a083137785cade5eb5))


### Bug Fixes

* **oscal:** deterministic OSCAL model write ([#553](https://github.com/defenseunicorns/lula/issues/553)) ([5493df1](https://github.com/defenseunicorns/lula/commit/5493df122b803d11542f29cfe80dfa4d5aaa10a8))


### Miscellaneous

* **deps:** update github/codeql-action action to v3.25.14 ([#557](https://github.com/defenseunicorns/lula/issues/557)) ([5bfd94f](https://github.com/defenseunicorns/lula/commit/5bfd94febc467e5a455ed32d97ce2e82e20409c2))
* **deps:** update github/codeql-action action to v3.25.15 ([#564](https://github.com/defenseunicorns/lula/issues/564)) ([60e128a](https://github.com/defenseunicorns/lula/commit/60e128a0a34ce8686c67e22ea2aebb61212b97fc))
* **deps:** update golang to version 1.22.5 ([#562](https://github.com/defenseunicorns/lula/issues/562)) ([97ff760](https://github.com/defenseunicorns/lula/commit/97ff7602f30f0709bd2ca16b74e53008607c3a61))
* **deps:** update module github.com/open-policy-agent/opa to v0.67.0 ([#561](https://github.com/defenseunicorns/lula/issues/561)) ([4378242](https://github.com/defenseunicorns/lula/commit/43782420b8b34362d03bcc965e00df2a850715c6))
* **docs:** fix simple demo command for evaluate file ([33fb97c](https://github.com/defenseunicorns/lula/commit/33fb97cccc9d4a589da65c03cc433b4f05c79d5d))
* **docs:** updated broken links ([#554](https://github.com/defenseunicorns/lula/issues/554)) ([8dd24b0](https://github.com/defenseunicorns/lula/commit/8dd24b083c86b12af8740fe788c4222f4c1c8718))
* **docs:** updated README for docs badge ([#558](https://github.com/defenseunicorns/lula/issues/558)) ([72fd3fc](https://github.com/defenseunicorns/lula/commit/72fd3fc8137477a4f10507481f8464eb5685b781))

## [0.4.3](https://github.com/defenseunicorns/lula/compare/v0.4.2...v0.4.3) (2024-07-19)


### Features

* **common:** json schema linting for common validation(s) ([#473](https://github.com/defenseunicorns/lula/issues/473)) ([23a45b6](https://github.com/defenseunicorns/lula/commit/23a45b696a3c24653ad2001dc4b883f40e9685c1))


### Bug Fixes

* **release:** add option to milestone for release process ([#535](https://github.com/defenseunicorns/lula/issues/535)) ([6fe64d8](https://github.com/defenseunicorns/lula/commit/6fe64d82ac4950214749b5f49a1ada12f43d193a))
* **test:** updated uuid in kyverno validation ([#539](https://github.com/defenseunicorns/lula/issues/539)) ([81446d9](https://github.com/defenseunicorns/lula/commit/81446d9441e1f062c57fa922e7d3cca833cbfd3e))


### Miscellaneous

* **deps:** update anchore/sbom-action action to v0.17.0 ([#541](https://github.com/defenseunicorns/lula/issues/541)) ([7c29fb7](https://github.com/defenseunicorns/lula/commit/7c29fb7dbbab163c648b4c04c89a1568206b8407))
* **deps:** update github/codeql-action action to v3.25.13 ([#507](https://github.com/defenseunicorns/lula/issues/507)) ([dc6cb88](https://github.com/defenseunicorns/lula/commit/dc6cb88eb8cda95c4f000988fc88e7ff1493d3cb))
* **deps:** update kubernetes packages to v0.30.3 ([#543](https://github.com/defenseunicorns/lula/issues/543)) ([1bdefce](https://github.com/defenseunicorns/lula/commit/1bdefce3f3e2af86f985f5b5e95d8d5f2c0c3c39))
* **docs:** initial docs structure/changes for feedback ([#524](https://github.com/defenseunicorns/lula/issues/524)) ([c276fdd](https://github.com/defenseunicorns/lula/commit/c276fdd3d390719e0a7825e0aabcdc50f0c33a0a))

## [0.4.2](https://github.com/defenseunicorns/lula/compare/v0.4.1...v0.4.2) (2024-07-10)


### Bug Fixes

* **evaluate:** set threshold on single result evaluation ([#519](https://github.com/defenseunicorns/lula/issues/519)) ([9424ec5](https://github.com/defenseunicorns/lula/commit/9424ec521f1ee1f4ddceb3350f22d4b3edea226d))
* **generate:** create annotation in remarks for how to reproduce the generation of a component ([#520](https://github.com/defenseunicorns/lula/issues/520)) ([6b59daf](https://github.com/defenseunicorns/lula/commit/6b59daffea89c82cd1b9b418f9b87cac81a3970e))
* **upgrade:** error handling for non-existent oscal ([#529](https://github.com/defenseunicorns/lula/issues/529)) ([58c03d5](https://github.com/defenseunicorns/lula/commit/58c03d528f05b42f98b67d7ba73d0ec86b3e5c9a))


### Miscellaneous

* **deps:** update actions/download-artifact action to v4.1.8 ([#522](https://github.com/defenseunicorns/lula/issues/522)) ([f628db8](https://github.com/defenseunicorns/lula/commit/f628db8a1df82d4357f289a77132839375b69df3))
* **deps:** update actions/setup-go action to v5.0.2 ([#530](https://github.com/defenseunicorns/lula/issues/530)) ([d6fa46e](https://github.com/defenseunicorns/lula/commit/d6fa46ef50d502ec1168282807ead1f4ea02c405))
* **deps:** update actions/setup-node action to v4.0.3 ([#526](https://github.com/defenseunicorns/lula/issues/526)) ([de146f7](https://github.com/defenseunicorns/lula/commit/de146f7283c94bce50bcf7f2492af8615dd1e523))
* **deps:** update actions/upload-artifact action to v4.3.4 ([#523](https://github.com/defenseunicorns/lula/issues/523)) ([1d2334b](https://github.com/defenseunicorns/lula/commit/1d2334b0ff676c32a5e3905db6e184d58872b5b7))
* **deps:** update anchore/sbom-action action to v0.16.1 ([#528](https://github.com/defenseunicorns/lula/issues/528)) ([ebdf05c](https://github.com/defenseunicorns/lula/commit/ebdf05caef149a0e21279e942169b96c4c883713))
* **website:** website removal ([#525](https://github.com/defenseunicorns/lula/issues/525)) ([575044c](https://github.com/defenseunicorns/lula/commit/575044c5c2b366ee160a2eb477a16a95192bc4e0))

## [0.4.1](https://github.com/defenseunicorns/lula/compare/v0.4.0...v0.4.1) (2024-06-29)


### Bug Fixes

* **release:** proper flag utilization ([#511](https://github.com/defenseunicorns/lula/issues/511)) ([ecefd9a](https://github.com/defenseunicorns/lula/commit/ecefd9ac094a590dd3d76212695605f26f4c3dcb))

## [0.4.0](https://github.com/defenseunicorns/lula/compare/v0.3.0...v0.4.0) (2024-06-28)


### ⚠ BREAKING CHANGES

* **validate:** #408 create resources in kubernetes domain ([#415](https://github.com/defenseunicorns/lula/issues/415))

### Features

* **tools:** add looping for lint ([#481](https://github.com/defenseunicorns/lula/issues/481)) ([0d69a45](https://github.com/defenseunicorns/lula/commit/0d69a45b6c8001a3923b9a66d5cfd7fcce3e1037))
* **validate:** [#408](https://github.com/defenseunicorns/lula/issues/408) create resources in kubernetes domain ([#415](https://github.com/defenseunicorns/lula/issues/415)) ([bd8d72b](https://github.com/defenseunicorns/lula/commit/bd8d72b1d3c31bbcf7fb9cdbc8fdfc38fccfcabb))


### Bug Fixes

* **commitlint:** pinned dependency issue by extracting into package(-… ([#454](https://github.com/defenseunicorns/lula/issues/454)) ([17ac8ca](https://github.com/defenseunicorns/lula/commit/17ac8ca76c081f2d385a8dca4bd35ced8cbbf70d))
* **evaluate:** add support for existing control-id becoming satisfied ([#498](https://github.com/defenseunicorns/lula/issues/498)) ([471e9c5](https://github.com/defenseunicorns/lula/commit/471e9c5a1b731b03adb64343954326b3ed4a4a3d))
* **evaluate:** establish threshold for assessment results result ([#457](https://github.com/defenseunicorns/lula/issues/457)) ([4571cb8](https://github.com/defenseunicorns/lula/commit/4571cb88a7fa1828387edbb7f35b44b65e09ad4f))
* **generate:** resolve parent flag options properly ([#442](https://github.com/defenseunicorns/lula/issues/442)) ([5850115](https://github.com/defenseunicorns/lula/commit/585011577ffcc9426a694296d631d29bc88b1f99))
* **oscal:** single model write operations support ([#502](https://github.com/defenseunicorns/lula/issues/502)) ([3646650](https://github.com/defenseunicorns/lula/commit/36466509aeff3aa53e08cb90454547537440b172))
* **validate:** fix related observations when empty ([#448](https://github.com/defenseunicorns/lula/issues/448)) ([f6f6993](https://github.com/defenseunicorns/lula/commit/f6f699312b1dfaa1190292143a039c105f5d5293))


### Miscellaneous

* **actions:** fix code scanning alerts ([#446](https://github.com/defenseunicorns/lula/issues/446)) ([aa568c7](https://github.com/defenseunicorns/lula/commit/aa568c76f9667195a0f106a4762f3ff315dc1666))
* **deps:** add tag to release please action ([#496](https://github.com/defenseunicorns/lula/issues/496)) ([3596491](https://github.com/defenseunicorns/lula/commit/35964912a7138721817b2ad136fe858fdb5dc8d3))
* **deps:** update actions/checkout action to v4.1.7 ([#479](https://github.com/defenseunicorns/lula/issues/479)) ([962fd2f](https://github.com/defenseunicorns/lula/commit/962fd2f6e694df92eb1a29735dd9750c13e3a0d9))
* **deps:** update dependency linkinator to v6.0.5 ([#458](https://github.com/defenseunicorns/lula/issues/458)) ([dfa1cbe](https://github.com/defenseunicorns/lula/commit/dfa1cbef30c40ccdba1b46915fd63aebc080f6e0))
* **deps:** update dependency markdownlint-cli to v0.41.0 ([#443](https://github.com/defenseunicorns/lula/issues/443)) ([27c0e94](https://github.com/defenseunicorns/lula/commit/27c0e947d3e9cf92a24687251acd786732e1a441))
* **deps:** update dependency prettier to v3.3.1 ([#459](https://github.com/defenseunicorns/lula/issues/459)) ([c999b78](https://github.com/defenseunicorns/lula/commit/c999b78c945e73ea1c5e93d269341e7c2c953e17))
* **deps:** update github/codeql-action action to v3.25.10 ([#476](https://github.com/defenseunicorns/lula/issues/476)) ([e2e74eb](https://github.com/defenseunicorns/lula/commit/e2e74eb473a6b2a315f908f21c4e2e56bd9e8b82))
* **deps:** update github/codeql-action action to v3.25.7 ([#452](https://github.com/defenseunicorns/lula/issues/452)) ([2583eea](https://github.com/defenseunicorns/lula/commit/2583eea6200f305c74ae214576188c424d39f749))
* **deps:** update github/codeql-action action to v3.25.8 ([#463](https://github.com/defenseunicorns/lula/issues/463)) ([0e7f844](https://github.com/defenseunicorns/lula/commit/0e7f844079826ac7fde2f5572cc5c48f377917a0))
* **deps:** update googleapis/release-please-action digest to 7987652 ([#472](https://github.com/defenseunicorns/lula/issues/472)) ([d50b034](https://github.com/defenseunicorns/lula/commit/d50b034802114ac3b1f3f26f50ec90640ee2f88a))
* **deps:** update goreleaser/goreleaser-action action to v6 ([#464](https://github.com/defenseunicorns/lula/issues/464)) ([e74b9d5](https://github.com/defenseunicorns/lula/commit/e74b9d53add588a4ddb83138b04e4e9ca263dae6))
* **deps:** update kubernetes packages to v0.30.2 ([#477](https://github.com/defenseunicorns/lula/issues/477)) ([679d2c8](https://github.com/defenseunicorns/lula/commit/679d2c8b684f07c7098465c14ffbbfd1c3e7d73b))
* **deps:** update module github.com/defenseunicorns/go-oscal to v0.4.3 ([#470](https://github.com/defenseunicorns/lula/issues/470)) ([5c78254](https://github.com/defenseunicorns/lula/commit/5c782545cdbd400f88a8e6e5e306a4a80207c0a2))
* **deps:** update module github.com/defenseunicorns/go-oscal to v0.5.0 ([#492](https://github.com/defenseunicorns/lula/issues/492)) ([c5d128f](https://github.com/defenseunicorns/lula/commit/c5d128ffaa831f4acdf6f6250732762eeb7d7397))
* **deps:** update module github.com/kyverno/kyverno-json to v0.0.3 ([#453](https://github.com/defenseunicorns/lula/issues/453)) ([1dc96e8](https://github.com/defenseunicorns/lula/commit/1dc96e832aa29a1d126a2b4f4d620aa006c8fec1))
* **deps:** update module github.com/open-policy-agent/opa to v0.65.0 ([#451](https://github.com/defenseunicorns/lula/issues/451)) ([7867a3c](https://github.com/defenseunicorns/lula/commit/7867a3c36826eaa264106434f993f37597ffd8e7))
* **deps:** update module github.com/open-policy-agent/opa to v0.66.0 ([#505](https://github.com/defenseunicorns/lula/issues/505)) ([7692e33](https://github.com/defenseunicorns/lula/commit/7692e3393bfa1ce26a3aca3acea125cf4928e61e))
* **deps:** update module github.com/spf13/cobra to v1.8.1 ([#485](https://github.com/defenseunicorns/lula/issues/485)) ([aaeba70](https://github.com/defenseunicorns/lula/commit/aaeba704b966ec65725f606647304482a77b9212))
* **docs:** cleanup unused readme conflicting in docs build ([#489](https://github.com/defenseunicorns/lula/issues/489)) ([19a3f61](https://github.com/defenseunicorns/lula/commit/19a3f61738a3812fa859aefcc86fdc53fa7c8246))
* **docs:** re-organize docs for docs website consumption ([#495](https://github.com/defenseunicorns/lula/issues/495)) ([24c24f0](https://github.com/defenseunicorns/lula/commit/24c24f0c8b5f40cd436c77e665fc068012ae41dc))
* fix documentation links ([#487](https://github.com/defenseunicorns/lula/issues/487)) ([4f96ec5](https://github.com/defenseunicorns/lula/commit/4f96ec5a8ebcbd9a7eb921a5e96db396c4967a02))

## [0.3.0](https://github.com/defenseunicorns/lula/compare/v0.2.1...v0.3.0) (2024-05-24)


### ⚠ BREAKING CHANGES

* #388 update types to use pointers ([#410](https://github.com/defenseunicorns/lula/issues/410))
* #367 compiling external/remote validations ([#384](https://github.com/defenseunicorns/lula/issues/384))

### refactor

* [#388](https://github.com/defenseunicorns/lula/issues/388) update types to use pointers ([#410](https://github.com/defenseunicorns/lula/issues/410)) ([9c51d56](https://github.com/defenseunicorns/lula/commit/9c51d5693565022f353c2739c97fac2686d78ce4))


### Features

* [#367](https://github.com/defenseunicorns/lula/issues/367) compiling external/remote validations ([#384](https://github.com/defenseunicorns/lula/issues/384)) ([8bb42b0](https://github.com/defenseunicorns/lula/commit/8bb42b02f6da2670830f11a1d2e1e5367c2b7d09))
* **oscal:** merge on write consolidation ([#407](https://github.com/defenseunicorns/lula/pull/407)) ([ef2f9f5](https://github.com/defenseunicorns/lula/commit/ef2f9f536ac8809785ca03325ae56575bbe0361c))
* **compose:** add ability to pull and compose import component defs ([#406](https://github.com/defenseunicorns/lula/pull/406)) ([ddf919a](https://github.com/defenseunicorns/lula/commit/ddf919a43995703f782f8667eb7305363b95b3cd))
* **generate:** add generate command and initial component generation ([#401](https://github.com/defenseunicorns/lula/issues/401)) ([918299a](https://github.com/defenseunicorns/lula/commit/918299ad397363e0d3580cd15b92ddf09044ce05))
* **dev:** added Observation logging to dev validate ([#396](https://github.com/defenseunicorns/lula/pull/396)) ([c32027e](https://github.com/defenseunicorns/lula/commit/c32027eafbf65e1cf69b3a72acac3a51c4f35656))
* **dev:** dev validate with optional resources file input ([#394](https://github.com/defenseunicorns/lula/pull/394)) ([f034a97](https://github.com/defenseunicorns/lula/commit/f034a976d20d10fe5ec660433e4554a02f76158c))
* **validate:** validation store/cache ([#373](https://github.com/defenseunicorns/lula/issues/373)) ([751982f](https://github.com/defenseunicorns/lula/commit/751982f5d4eee60a6412eed5e554c86a683ecb7a))


### Bug Fixes

* **deps:** consolidate use of goyaml pkg ([#422](https://github.com/defenseunicorns/lula/issues/422)) ([d1abbcc](https://github.com/defenseunicorns/lula/commit/d1abbcc052fd1f1ff2e57265e54a1b005ec66641))
* **deps:** controller runtime ([#379](https://github.com/defenseunicorns/lula/issues/379)) ([7d3aec3](https://github.com/defenseunicorns/lula/commit/7d3aec3e7e94652d524d1e40d62c61736ca1e12e))
* **dev:** updated result condition to match satisfaction logic ([#400](https://github.com/defenseunicorns/lula/issues/400)) ([5feda9d](https://github.com/defenseunicorns/lula/commit/5feda9dde93b270e9d2cebc5ee40ec21ab4b1c4c))
* **validate:** validation errors mapped to observations ([#405](https://github.com/defenseunicorns/lula/pull/405)) ([39e5ebd](https://github.com/defenseunicorns/lula/commit/39e5ebd45d4a9cfc0918899ed647192d8bcf0952))
* **validate:** fix order of assessment-results results ([#437](https://github.com/defenseunicorns/lula/issues/437)) ([a8db208](https://github.com/defenseunicorns/lula/commit/a8db20862f9f1bf7c269cd75839823f58b7c9541))


### Miscellaneous

* **actions:** [#420](https://github.com/defenseunicorns/lula/issues/420) update release process with release-please ([#421](https://github.com/defenseunicorns/lula/issues/421)) ([a372df0](https://github.com/defenseunicorns/lula/commit/a372df0e0316067a3adb62d02c95433d37930ec5))
* **deps:** bump golang.org/x/net from 0.22.0 to 0.23.0 ([#378](https://github.com/defenseunicorns/lula/issues/378)) ([8088bd0](https://github.com/defenseunicorns/lula/commit/8088bd0d38c89768a7ee4eae7e12edea3ce4f35e))
* **deps:** Update actions/checkout action to v4.1.3 ([#382](https://github.com/defenseunicorns/lula/issues/382)) ([08eed39](https://github.com/defenseunicorns/lula/commit/08eed39078dc38e79b5b5483b2cdd0770a00594d))
* **deps:** Update actions/download-artifact action to v4.1.6 ([#376](https://github.com/defenseunicorns/lula/issues/376)) ([2982b36](https://github.com/defenseunicorns/lula/commit/2982b3659d00671bdd6eac0a29ad25aca0e7da30))
* **deps:** Update actions/download-artifact action to v4.1.7 ([#387](https://github.com/defenseunicorns/lula/issues/387)) ([92064e6](https://github.com/defenseunicorns/lula/commit/92064e6184581c116d0fd4fd07521c852211ebcf))
* **deps:** Update actions/upload-artifact action to v4.3.2 ([#377](https://github.com/defenseunicorns/lula/issues/377)) ([f575f82](https://github.com/defenseunicorns/lula/commit/f575f82a3b00d2c5260b0887391638faa427b8ee))
* **deps:** Update actions/upload-artifact action to v4.3.3 ([#383](https://github.com/defenseunicorns/lula/issues/383)) ([26f1f32](https://github.com/defenseunicorns/lula/commit/26f1f32d54e1fe64f528ca3441cbd39055dc8ee2))
* **deps:** update anchore/sbom-action action to v0.16.0 ([#426](https://github.com/defenseunicorns/lula/issues/426)) ([a2063a5](https://github.com/defenseunicorns/lula/commit/a2063a5d47e9d816901cda1f908c6e6d68b53442))
* **deps:** update github/codeql-action action to v3.25.6 ([#425](https://github.com/defenseunicorns/lula/issues/425)) ([9ef1703](https://github.com/defenseunicorns/lula/commit/9ef1703f0180f05bec97b1fa5f894fae5d4627f6))
* **deps:** update golang to version 1.22.3 ([#423](https://github.com/defenseunicorns/lula/issues/423)) ([aa8cab7](https://github.com/defenseunicorns/lula/commit/aa8cab732053ed588ec7c6d83895e3d6f0ecf7f4))
* **deps:** update kubernetes packages to v0.30.1 ([#417](https://github.com/defenseunicorns/lula/issues/417)) ([e47a04d](https://github.com/defenseunicorns/lula/commit/e47a04d4df9a9f3f14a157716120cf5f12714d5c))
* **deps:** Update module github.com/defenseunicorns/go-oscal to v0.3.2 ([#380](https://github.com/defenseunicorns/lula/issues/380)) ([03aa969](https://github.com/defenseunicorns/lula/commit/03aa969ff102111325c0b45b91e7c3543c15cf16))
* **deps:** update module github.com/defenseunicorns/go-oscal to v0.4.0 ([#429](https://github.com/defenseunicorns/lula/issues/429)) ([4ff7775](https://github.com/defenseunicorns/lula/commit/4ff7775890113407c60240ca21382dfca0eb102c))
* **deps:** update module github.com/defenseunicorns/go-oscal to v0.4.1 ([#435](https://github.com/defenseunicorns/lula/issues/435)) ([4570658](https://github.com/defenseunicorns/lula/commit/4570658fcc3d20d7c4d118e89626b5d81267af91))
* **deps:** update module github.com/hashicorp/go-version to v1.7.0 ([#438](https://github.com/defenseunicorns/lula/issues/438)) ([4f6de9b](https://github.com/defenseunicorns/lula/commit/4f6de9b0cff92530a6d0e34697bfa694d4e86f33))
* **deps:** update module sigs.k8s.io/e2e-framework to v0.4.0 ([#419](https://github.com/defenseunicorns/lula/issues/419)) ([890a7d8](https://github.com/defenseunicorns/lula/commit/890a7d8d7a5da123b6dff23b75e1b390ff7ca716))
* **renovate:** update config to handle conventional commit titles ([#428](https://github.com/defenseunicorns/lula/issues/428)) ([5f4139a](https://github.com/defenseunicorns/lula/commit/5f4139a3b6df6fd5ba4c1ee7f4e04dd50f23be1f))
