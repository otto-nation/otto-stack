# Changelog

## [1.2.1](https://github.com/otto-nation/otto-stack/compare/v1.2.0...v1.2.1) (2026-03-15)


### Bug Fixes

* **ci:** simplify release config, fix brew tap, and improve UX ([#53](https://github.com/otto-nation/otto-stack/issues/53)) ([cece57c](https://github.com/otto-nation/otto-stack/commit/cece57cdbad9b59e4d18ad8db2e2d2f4d6a81868))

## [1.2.0](https://github.com/otto-nation/otto-stack/compare/v1.1.0...v1.2.0) (2026-03-13)


### Features

* add --force flag support for sharing non-shareable services ([73fb2e3](https://github.com/otto-nation/otto-stack/commit/73fb2e3fbb00c6a0f425cfc965229a646f5068dd))
* add complete CLI implementation for otto-stack ([a62742a](https://github.com/otto-nation/otto-stack/commit/a62742af6247a62595406c1aa8266e6f577f5630))
* add comprehensive CLI installation and build automation ([c13f7aa](https://github.com/otto-nation/otto-stack/commit/c13f7aad923a17e7d9067e45b75367facce25dc2))
* add configurable AI agent support for git automation ([ece9966](https://github.com/otto-nation/otto-stack/commit/ece99661b50cf80710366bf72fbeed6be2043c9d))
* add orphan detection and safety features for shared containers ([dac3f02](https://github.com/otto-nation/otto-stack/commit/dac3f02c020639bfc59e9584517bbd57da9ae10a))
* add shareable field to service YAML schema ([feff381](https://github.com/otto-nation/otto-stack/commit/feff3810ad6aa3f7f9850fe958ab7f686ab00a20))
* add Shareable field to ServiceConfig and enforce it ([bbebb0e](https://github.com/otto-nation/otto-stack/commit/bbebb0ec772cde3e3447ed6aad338d937faa73a2))
* add shared container flags to down command and fix naming ([3c5232b](https://github.com/otto-nation/otto-stack/commit/3c5232bb0878cf961dce391de04d57f15742c794))
* add shared container flags to init command ([12ac1c1](https://github.com/otto-nation/otto-stack/commit/12ac1c15116eeb7af9cd992976099ef2c75c7968))
* add shared context support to all operation handlers ([8519761](https://github.com/otto-nation/otto-stack/commit/8519761b786e3d1d55c693bf7025abb9ebc81f28))
* add sharing configuration to project config ([3586746](https://github.com/otto-nation/otto-stack/commit/3586746ff96d71bf87f2e5b08f6dda5f1e6611c4))
* add sharing metrics and visibility to status command ([b9bf716](https://github.com/otto-nation/otto-stack/commit/b9bf71659a83c663c75422f6b5040f67811138c3))
* add sharing policy validation during config load ([2acd3ee](https://github.com/otto-nation/otto-stack/commit/2acd3ee8e2a3223dedb32ddb3b0212f7775090af))
* auto-generate flag validation tests from commands.yaml ([fa8b161](https://github.com/otto-nation/otto-stack/commit/fa8b16121d11c5a06b284ea08af77d757ff61c33))
* enhance orphan detection with filesystem and Docker checks ([1d56cea](https://github.com/otto-nation/otto-stack/commit/1d56ceac2565adcd8c4d59871aeef838680958da))
* implement container naming and lifecycle strategy ([9e27b6c](https://github.com/otto-nation/otto-stack/commit/9e27b6cd1cbbc2bd06698a95d51c5ede93c9750a))
* implement context detection system ([e5f7e69](https://github.com/otto-nation/otto-stack/commit/e5f7e691bba7b683a934cf8b494b29dc0efecee5))
* implement context-aware lifecycle commands (up/down/status) ([de7ebd5](https://github.com/otto-nation/otto-stack/commit/de7ebd5a88fad303021137483b51a42f33d224d0))
* implement docker-compose pass-through flags and remove unimplemented flags ([f807d50](https://github.com/otto-nation/otto-stack/commit/f807d506e56d7e17e5cc9b6cf4be096671a8c8d1))
* implement registry reconciliation with Docker state ([4d21b85](https://github.com/otto-nation/otto-stack/commit/4d21b85462496c5e44141f1fc59d467964f8afca))
* implement shared container registry ([a365953](https://github.com/otto-nation/otto-stack/commit/a36595362e94c2e6c9af88dafe26d61d8bc396a7))
* integrate validation functions and refactor service layer ([3e573f3](https://github.com/otto-nation/otto-stack/commit/3e573f3271f580af932444ad08e4740d7680eb69))
* **registry:** shared registry correctness, enriched status, and CLI polish ([9a7a0c5](https://github.com/otto-nation/otto-stack/commit/9a7a0c523301e6e321f03aa42ecaba7cd3cc02ed))
* shared containers, middleware chain, registry reconciliation, and lifecycle fixes ([#44](https://github.com/otto-nation/otto-stack/issues/44)) ([3bd5417](https://github.com/otto-nation/otto-stack/commit/3bd54172f8240d3374af74faab697cf2118fdab2))
* **shared:** start/stop via compose; files moved to generated/ ([85888bc](https://github.com/otto-nation/otto-stack/commit/85888bc514e0afb95cc711d9c3038ad9cb95bc3b))


### Bug Fixes

* **ci:** enable homebrew tap and clean up release config generation ([#51](https://github.com/otto-nation/otto-stack/issues/51)) ([6550b45](https://github.com/otto-nation/otto-stack/commit/6550b4587790f6833cec4daaa41e8d918aadcacf))
* implement file locking for registry to prevent corruption ([479e776](https://github.com/otto-nation/otto-stack/commit/479e77650c0bc2eeae2b847cf5104a409cfc6925))
* remove docs-lastmod-update from pre-push hook ([5cdbe9b](https://github.com/otto-nation/otto-stack/commit/5cdbe9bb4bf9617ffa25dbf236383cad4e5dac19))
* update OpenTelemetry SDK to v1.40.0 and allow Dependabot commits ([3593ba8](https://github.com/otto-nation/otto-stack/commit/3593ba82d155cba02c86538fb1136389e12f8d8d))


### Dependencies

* bump the npm-dependencies group across 1 directory with 6 updates ([#11](https://github.com/otto-nation/otto-stack/issues/11)) ([a877852](https://github.com/otto-nation/otto-stack/commit/a877852a973d6c5c7ed4d2fc3d09812dcaf1233c))
* bump the npm-dependencies group in /docs-site with 16 updates ([c7078ee](https://github.com/otto-nation/otto-stack/commit/c7078eed7e654cf9c323cd6f94af3bd931515b84))
* bump the npm-dependencies group in /docs-site with 5 updates ([d6a891d](https://github.com/otto-nation/otto-stack/commit/d6a891d5cfc3ea6d95cbf2fdd1c1ef61687254bd))


### Reverts

* **main:** release otto-stack 1.1.0 ([#45](https://github.com/otto-nation/otto-stack/issues/45)) ([#49](https://github.com/otto-nation/otto-stack/issues/49)) ([ee58c24](https://github.com/otto-nation/otto-stack/commit/ee58c24b17f956d8a386709879b85698da4c88da))

## [1.1.0](https://github.com/otto-nation/otto-stack/compare/otto-stack-v1.0.0...otto-stack-v1.1.0) (2026-03-13)


### Features

* shared containers, middleware chain, registry reconciliation, and lifecycle fixes ([#44](https://github.com/otto-nation/otto-stack/issues/44)) ([3bd5417](https://github.com/otto-nation/otto-stack/commit/3bd54172f8240d3374af74faab697cf2118fdab2))


### Dependencies

* bump the npm-dependencies group across 1 directory with 6 updates ([#11](https://github.com/otto-nation/otto-stack/issues/11)) ([a877852](https://github.com/otto-nation/otto-stack/commit/a877852a973d6c5c7ed4d2fc3d09812dcaf1233c))


### Reverts

* **main:** release otto-stack 1.1.0 ([#45](https://github.com/otto-nation/otto-stack/issues/45)) ([#49](https://github.com/otto-nation/otto-stack/issues/49)) ([ee58c24](https://github.com/otto-nation/otto-stack/commit/ee58c24b17f956d8a386709879b85698da4c88da))

## 1.0.0 (2026-02-26)


### Features

* add --force flag support for sharing non-shareable services ([73fb2e3](https://github.com/otto-nation/otto-stack/commit/73fb2e3fbb00c6a0f425cfc965229a646f5068dd))
* add complete CLI implementation for otto-stack ([a62742a](https://github.com/otto-nation/otto-stack/commit/a62742af6247a62595406c1aa8266e6f577f5630))
* add comprehensive CLI installation and build automation ([c13f7aa](https://github.com/otto-nation/otto-stack/commit/c13f7aad923a17e7d9067e45b75367facce25dc2))
* add configurable AI agent support for git automation ([ece9966](https://github.com/otto-nation/otto-stack/commit/ece99661b50cf80710366bf72fbeed6be2043c9d))
* add orphan detection and safety features for shared containers ([dac3f02](https://github.com/otto-nation/otto-stack/commit/dac3f02c020639bfc59e9584517bbd57da9ae10a))
* add shareable field to service YAML schema ([feff381](https://github.com/otto-nation/otto-stack/commit/feff3810ad6aa3f7f9850fe958ab7f686ab00a20))
* add Shareable field to ServiceConfig and enforce it ([bbebb0e](https://github.com/otto-nation/otto-stack/commit/bbebb0ec772cde3e3447ed6aad338d937faa73a2))
* add shared container flags to down command and fix naming ([3c5232b](https://github.com/otto-nation/otto-stack/commit/3c5232bb0878cf961dce391de04d57f15742c794))
* add shared container flags to init command ([12ac1c1](https://github.com/otto-nation/otto-stack/commit/12ac1c15116eeb7af9cd992976099ef2c75c7968))
* add shared context support to all operation handlers ([8519761](https://github.com/otto-nation/otto-stack/commit/8519761b786e3d1d55c693bf7025abb9ebc81f28))
* add sharing configuration to project config ([3586746](https://github.com/otto-nation/otto-stack/commit/3586746ff96d71bf87f2e5b08f6dda5f1e6611c4))
* add sharing metrics and visibility to status command ([b9bf716](https://github.com/otto-nation/otto-stack/commit/b9bf71659a83c663c75422f6b5040f67811138c3))
* add sharing policy validation during config load ([2acd3ee](https://github.com/otto-nation/otto-stack/commit/2acd3ee8e2a3223dedb32ddb3b0212f7775090af))
* auto-generate flag validation tests from commands.yaml ([fa8b161](https://github.com/otto-nation/otto-stack/commit/fa8b16121d11c5a06b284ea08af77d757ff61c33))
* enhance orphan detection with filesystem and Docker checks ([1d56cea](https://github.com/otto-nation/otto-stack/commit/1d56ceac2565adcd8c4d59871aeef838680958da))
* implement container naming and lifecycle strategy ([9e27b6c](https://github.com/otto-nation/otto-stack/commit/9e27b6cd1cbbc2bd06698a95d51c5ede93c9750a))
* implement context detection system ([e5f7e69](https://github.com/otto-nation/otto-stack/commit/e5f7e691bba7b683a934cf8b494b29dc0efecee5))
* implement context-aware lifecycle commands (up/down/status) ([de7ebd5](https://github.com/otto-nation/otto-stack/commit/de7ebd5a88fad303021137483b51a42f33d224d0))
* implement docker-compose pass-through flags and remove unimplemented flags ([f807d50](https://github.com/otto-nation/otto-stack/commit/f807d506e56d7e17e5cc9b6cf4be096671a8c8d1))
* implement registry reconciliation with Docker state ([4d21b85](https://github.com/otto-nation/otto-stack/commit/4d21b85462496c5e44141f1fc59d467964f8afca))
* implement shared container registry ([a365953](https://github.com/otto-nation/otto-stack/commit/a36595362e94c2e6c9af88dafe26d61d8bc396a7))
* integrate validation functions and refactor service layer ([3e573f3](https://github.com/otto-nation/otto-stack/commit/3e573f3271f580af932444ad08e4740d7680eb69))
* **registry:** shared registry correctness, enriched status, and CLI polish ([9a7a0c5](https://github.com/otto-nation/otto-stack/commit/9a7a0c523301e6e321f03aa42ecaba7cd3cc02ed))
* **shared:** start/stop via compose; files moved to generated/ ([85888bc](https://github.com/otto-nation/otto-stack/commit/85888bc514e0afb95cc711d9c3038ad9cb95bc3b))


### Bug Fixes

* implement file locking for registry to prevent corruption ([479e776](https://github.com/otto-nation/otto-stack/commit/479e77650c0bc2eeae2b847cf5104a409cfc6925))
* remove docs-lastmod-update from pre-push hook ([5cdbe9b](https://github.com/otto-nation/otto-stack/commit/5cdbe9bb4bf9617ffa25dbf236383cad4e5dac19))
* update OpenTelemetry SDK to v1.40.0 and allow Dependabot commits ([3593ba8](https://github.com/otto-nation/otto-stack/commit/3593ba82d155cba02c86538fb1136389e12f8d8d))


### Dependencies

* bump the npm-dependencies group in /docs-site with 16 updates ([c7078ee](https://github.com/otto-nation/otto-stack/commit/c7078eed7e654cf9c323cd6f94af3bd931515b84))
* bump the npm-dependencies group in /docs-site with 5 updates ([d6a891d](https://github.com/otto-nation/otto-stack/commit/d6a891d5cfc3ea6d95cbf2fdd1c1ef61687254bd))
