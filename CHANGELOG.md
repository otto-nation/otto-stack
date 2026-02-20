# Changelog

## 1.0.0 (2026-02-20)


### Features

* add --force flag support for sharing non-shareable services ([#110](https://github.com/otto-nation/otto-stack/issues/110)) ([80fc36d](https://github.com/otto-nation/otto-stack/commit/80fc36d392dc6626d4784dc9eaf5d23f0144c1d7)), closes [#101](https://github.com/otto-nation/otto-stack/issues/101)
* add Close method to logger for proper file handle cleanup ([#5](https://github.com/otto-nation/otto-stack/issues/5)) ([0de0dff](https://github.com/otto-nation/otto-stack/commit/0de0dff42f9c6b8dd5aa3e755edfa338a0e45738))
* Add complete CLI implementation for otto-stack ([#1](https://github.com/otto-nation/otto-stack/issues/1)) ([456f050](https://github.com/otto-nation/otto-stack/commit/456f050739196b6eb55c961d7a702fc5770ccea1))
* add comprehensive CLI installation and build automation ([#12](https://github.com/otto-nation/otto-stack/issues/12)) ([d7fee41](https://github.com/otto-nation/otto-stack/commit/d7fee41e6b417b44549cee6746b42cbad2806866))
* add configurable AI agent support for git automation ([#23](https://github.com/otto-nation/otto-stack/issues/23)) ([b431b63](https://github.com/otto-nation/otto-stack/commit/b431b6353b84fb47840335cfb4ec8ffe1eadd304))
* add deadcode analysis and consolidate UI constants ([#85](https://github.com/otto-nation/otto-stack/issues/85)) ([2487cbb](https://github.com/otto-nation/otto-stack/commit/2487cbbe83c5d0b5c1d90c13f31b0ed39ee1cd98))
* add orphan detection and safety features for shared containers ([#80](https://github.com/otto-nation/otto-stack/issues/80)) ([d075e62](https://github.com/otto-nation/otto-stack/commit/d075e625392260a1857821e934d83f76c84a2ec5))
* add schema-driven config generation and flexible service output formats ([#25](https://github.com/otto-nation/otto-stack/issues/25)) ([3289b74](https://github.com/otto-nation/otto-stack/commit/3289b741975ce12627915111250aab0712c72f06))
* add shareable field to service YAML schema ([#63](https://github.com/otto-nation/otto-stack/issues/63)) ([441310b](https://github.com/otto-nation/otto-stack/commit/441310bfd346c94b6d70e383b65a34b0f8b9189d))
* Add Shareable field to ServiceConfig and enforce it ([#107](https://github.com/otto-nation/otto-stack/issues/107)) ([220259f](https://github.com/otto-nation/otto-stack/commit/220259f4f53bb4c84391bdde28598e5eef895b58)), closes [#99](https://github.com/otto-nation/otto-stack/issues/99)
* add shared container flags to down command and fix naming ([#118](https://github.com/otto-nation/otto-stack/issues/118)) ([cb598ae](https://github.com/otto-nation/otto-stack/commit/cb598ae4b80f3a00fb9364356e4a677fdb443bb0))
* add shared container flags to init command ([#69](https://github.com/otto-nation/otto-stack/issues/69)) ([d6b41cd](https://github.com/otto-nation/otto-stack/commit/d6b41cd0692b24c04411c2180717b792f9f17234))
* add shared context support to all operation handlers ([#116](https://github.com/otto-nation/otto-stack/issues/116)) ([0b0903b](https://github.com/otto-nation/otto-stack/commit/0b0903b3ff6d91a69c3d127f6c4f2d87ab430765)), closes [#104](https://github.com/otto-nation/otto-stack/issues/104)
* add sharing configuration to project config ([#66](https://github.com/otto-nation/otto-stack/issues/66)) ([7ce79c2](https://github.com/otto-nation/otto-stack/commit/7ce79c2c1170f6e37901fcaecfb94bda5041ae66))
* add sharing metrics and visibility to status command ([#117](https://github.com/otto-nation/otto-stack/issues/117)) ([1f85973](https://github.com/otto-nation/otto-stack/commit/1f859737a177701044ba3dee7bc4ad1e77d817d0))
* add sharing policy validation during config load ([#109](https://github.com/otto-nation/otto-stack/issues/109)) ([ca4f748](https://github.com/otto-nation/otto-stack/commit/ca4f748bed4789cf89681c71e66a2f1f831cabbb))
* auto-generate flag validation tests from commands.yaml ([#82](https://github.com/otto-nation/otto-stack/issues/82)) ([ce74b8a](https://github.com/otto-nation/otto-stack/commit/ce74b8a801993a7f9c2b53742ebade090cba91c6))
* consolidate CLI architecture and centralize branding constants ([#24](https://github.com/otto-nation/otto-stack/issues/24)) ([a4b14ea](https://github.com/otto-nation/otto-stack/commit/a4b14eabe39474a8223b23cc2878d3ea418267a7))
* enhance orphan detection with filesystem and Docker checks ([#115](https://github.com/otto-nation/otto-stack/issues/115)) ([669d0ea](https://github.com/otto-nation/otto-stack/commit/669d0eace7078de5fa4353f6b8411e11f184fdeb)), closes [#103](https://github.com/otto-nation/otto-stack/issues/103)
* implement container naming and lifecycle strategy ([#68](https://github.com/otto-nation/otto-stack/issues/68)) ([9749f69](https://github.com/otto-nation/otto-stack/commit/9749f69be23a749a8250dbfe07e6ffd95400caac))
* implement context detection system ([#64](https://github.com/otto-nation/otto-stack/issues/64)) ([4b27964](https://github.com/otto-nation/otto-stack/commit/4b279645c4d0fa9b0d408202c67ae0decc0eceea))
* implement context-aware lifecycle commands (up/down/status) ([#71](https://github.com/otto-nation/otto-stack/issues/71)) ([ed0cebe](https://github.com/otto-nation/otto-stack/commit/ed0cebed7f39ea3f6c5b8042d2b6928a1ec18a6f))
* implement docker-compose pass-through flags and remove unimplemented flags ([#81](https://github.com/otto-nation/otto-stack/issues/81)) ([2a4f07b](https://github.com/otto-nation/otto-stack/commit/2a4f07bbf378465f21d020c64e024330638145ea)), closes [#73](https://github.com/otto-nation/otto-stack/issues/73)
* implement registry reconciliation with Docker state ([#108](https://github.com/otto-nation/otto-stack/issues/108)) ([cc826b8](https://github.com/otto-nation/otto-stack/commit/cc826b8083311d0af91ff4caf39f775fef01382c)), closes [#100](https://github.com/otto-nation/otto-stack/issues/100)
* implement shared container registry ([#67](https://github.com/otto-nation/otto-stack/issues/67)) ([6fd4196](https://github.com/otto-nation/otto-stack/commit/6fd41962b84f37d8780179e1e5db4caf824c832a))
* integrate validation functions and refactor service layer ([#89](https://github.com/otto-nation/otto-stack/issues/89)) ([5002d4c](https://github.com/otto-nation/otto-stack/commit/5002d4cd8a08a276abc11c58186b5209043c29eb))
* refactor to Go-idiomatic interface-based type discrimination ([#119](https://github.com/otto-nation/otto-stack/issues/119)) ([32b40f0](https://github.com/otto-nation/otto-stack/commit/32b40f0a4688feaa484a2557afe5874a2092b048)), closes [#112](https://github.com/otto-nation/otto-stack/issues/112)


### Bug Fixes

* implement file locking for registry to prevent corruption ([#111](https://github.com/otto-nation/otto-stack/issues/111)) ([fe2d271](https://github.com/otto-nation/otto-stack/commit/fe2d271fb3cac70acb3018f387eca4d0f40746ef)), closes [#105](https://github.com/otto-nation/otto-stack/issues/105)
* improve cross-platform command detection and Windows compatibility ([#4](https://github.com/otto-nation/otto-stack/issues/4)) ([557f5ad](https://github.com/otto-nation/otto-stack/commit/557f5adcbfc1cc858e39513ca28675fd2ad422f4))
* remove docs-lastmod-update from pre-push hook ([#65](https://github.com/otto-nation/otto-stack/issues/65)) ([51d62d9](https://github.com/otto-nation/otto-stack/commit/51d62d96be4e0a725aca10fb70886185667ce011))
* remove unreliable LoadError tests ([856bb1d](https://github.com/otto-nation/otto-stack/commit/856bb1dce6c545f081799d74349f5b65600a3759))
* resolve YAML syntax error in CI stability monitor ([be4d719](https://github.com/otto-nation/otto-stack/commit/be4d71950e46f8446852f063206504a7b4be3a2d))
* skip file-as-directory tests on Windows ([9e3d239](https://github.com/otto-nation/otto-stack/commit/9e3d239bcb5fb6b23ad311c487e3024e6da10d87))
* update OpenTelemetry SDK to v1.40.0 and allow Dependabot commits ([#5](https://github.com/otto-nation/otto-stack/issues/5)) ([590ca4b](https://github.com/otto-nation/otto-stack/commit/590ca4b9dbf6e6036a0fb3264174406b303cf03d))
* use binary data for invalid YAML tests ([239b6c2](https://github.com/otto-nation/otto-stack/commit/239b6c245a37edc3f20ebd9880e903bc404cb1b4))
* use correct YAML field name in error tests ([fc51d44](https://github.com/otto-nation/otto-stack/commit/fc51d44f2047fd6ae5775c8879a0a26eceb126f3))
* use cross-platform error conditions in registry tests ([4c8cdc3](https://github.com/otto-nation/otto-stack/commit/4c8cdc3df95e8ed81486e6653e7e62aaf32dc738))
* use more robust invalid YAML for error tests ([f4edf6b](https://github.com/otto-nation/otto-stack/commit/f4edf6b4b92a1fd12558dfa0095565c21113c0f9))
* use tab character for invalid YAML test ([fc2eff8](https://github.com/otto-nation/otto-stack/commit/fc2eff8227649b078a4c6ae70663793ede04af6e))
* use type mismatch for YAML error tests ([b3448ef](https://github.com/otto-nation/otto-stack/commit/b3448ef5e0a07914274d62b5780bb8303c1dd4c9))


### Dependencies

* bump the npm-dependencies group in /docs-site with 16 updates ([#135](https://github.com/otto-nation/otto-stack/issues/135)) ([77c9a74](https://github.com/otto-nation/otto-stack/commit/77c9a74be3c0ea3aa113d019450c75321d9c22ec))
* **deps-dev:** bump markdown-link-check from 3.13.7 to 3.14.1 in /docs-site ([#6](https://github.com/otto-nation/otto-stack/issues/6)) ([4fec5d3](https://github.com/otto-nation/otto-stack/commit/4fec5d3cf834003f5d476366ec48bf9240f225a2))
* **deps-dev:** bump vite from 7.1.9 to 7.1.12 in /docs-site ([#20](https://github.com/otto-nation/otto-stack/issues/20)) ([c8a0a1a](https://github.com/otto-nation/otto-stack/commit/c8a0a1a2182d2479747da6cce6dd8d06d93b3ec9))
* **deps:** bump golang.org/x/text from 0.32.0 to 0.33.0 ([#37](https://github.com/otto-nation/otto-stack/issues/37)) ([b386205](https://github.com/otto-nation/otto-stack/commit/b386205bd282c4c97bffe8415f79eebe4845a5a8))
