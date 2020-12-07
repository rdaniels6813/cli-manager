## [1.2.6](https://github.com/rdaniels6813/cli-manager/compare/v1.2.5...v1.2.6) (2020-12-07)


### Bug Fixes

* get next version from semantic-release run ([710e5f7](https://github.com/rdaniels6813/cli-manager/commit/710e5f7f211b0f0965d932614c7fef779cee2d7f))

## [1.2.5](https://github.com/rdaniels6813/cli-manager/compare/v1.2.4...v1.2.5) (2020-12-07)


### Bug Fixes

* include version information with the CLI ([be21e1c](https://github.com/rdaniels6813/cli-manager/commit/be21e1ce96f3dad0e9e93f78672202ac2c143a34))

## [1.2.4](https://github.com/rdaniels6813/cli-manager/compare/v1.2.3...v1.2.4) (2020-11-20)


### Bug Fixes

* allow running commands to capture sigint ([a49c9cd](https://github.com/rdaniels6813/cli-manager/commit/a49c9cdb01365647885baa9251485037881d8631))

## [1.2.3](https://github.com/rdaniels6813/cli-manager/compare/v1.2.2...v1.2.3) (2020-10-29)


### Bug Fixes

* more bugfixes around engines.node ([b4d0263](https://github.com/rdaniels6813/cli-manager/commit/b4d02630b5a76f774639f7d940c477adb927f451))

## [1.2.2](https://github.com/rdaniels6813/cli-manager/compare/v1.2.1...v1.2.2) (2020-10-29)


### Bug Fixes

* use better version parser for engines.node range ([bd15dac](https://github.com/rdaniels6813/cli-manager/commit/bd15dac823a578a166eaf9c011bd8ecad403e3f6))

## [1.2.1](https://github.com/rdaniels6813/cli-manager/compare/v1.2.0...v1.2.1) (2020-09-30)


### Bug Fixes

* fix a bug with bin map being map[string]interface{} ([8c15e16](https://github.com/rdaniels6813/cli-manager/commit/8c15e1647d68d73afc8cc51bf8988466caed6820))

# [1.2.0](https://github.com/rdaniels6813/cli-manager/compare/v1.1.0...v1.2.0) (2020-09-30)


### Features

* shorten direct paths to bin to just bin name ([be35132](https://github.com/rdaniels6813/cli-manager/commit/be35132661e1891bdba190072b750e228c687ab2))

# [1.1.0](https://github.com/rdaniels6813/cli-manager/compare/v1.0.2...v1.1.0) (2020-09-30)


### Bug Fixes

* add github credentials to semantic-release ([48f9693](https://github.com/rdaniels6813/cli-manager/commit/48f9693be1717cbecf056af8cd727f3500ef9e22))
* remove circleci and fix github action build ([9f643a2](https://github.com/rdaniels6813/cli-manager/commit/9f643a29281ee859a4cb5fcfae17fbd5fc4fa609))
* updates to CI ([39f91c4](https://github.com/rdaniels6813/cli-manager/commit/39f91c4bf983875f137056382180e2738a6d2a97))


### Features

* upgrade go version to 1.15 ([b634a7f](https://github.com/rdaniels6813/cli-manager/commit/b634a7fa6a10b8cdd98e4d0ef7617f7119e2803d))

## [1.0.2](https://github.com/rdaniels6813/cli-manager/compare/v1.0.1...v1.0.2) (2019-08-22)


### Bug Fixes

* revert ([07d3d4c](https://github.com/rdaniels6813/cli-manager/commit/07d3d4c))

## [1.0.1](https://github.com/rdaniels6813/cli-manager/compare/v1.0.0...v1.0.1) (2019-08-22)


### Bug Fixes

* update tests in make file ([d604c10](https://github.com/rdaniels6813/cli-manager/commit/d604c10))

# 1.0.0 (2019-08-19)


### Bug Fixes

* add newlines to completions & aliases snippets ([31a83f6](https://github.com/rdaniels6813/cli-manager/commit/31a83f6))
* better zsh detection & use scanner to read profile ([c826d59](https://github.com/rdaniels6813/cli-manager/commit/c826d59))
* extend support for aliases & completions ([8d4d757](https://github.com/rdaniels6813/cli-manager/commit/8d4d757))
* go mod tidy ([bbd1d18](https://github.com/rdaniels6813/cli-manager/commit/bbd1d18))
* use the correct node version on unix OS ([d2c8695](https://github.com/rdaniels6813/cli-manager/commit/d2c8695))


### Features

* add uninstall command with support for package name & install name ([4285517](https://github.com/rdaniels6813/cli-manager/commit/4285517))
* configure completion & aliases for installed CLIS ([2618bb3](https://github.com/rdaniels6813/cli-manager/commit/2618bb3))
