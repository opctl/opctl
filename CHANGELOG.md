# Change Log

All notable changes to this project will be documented in this file.

## 0.1.2 - 2016-06-22
### Fixed
- [Missleading `variable is not set` message on op finish](https://github.com/opctl/engine/issues/5)
- [Engine not observing exitcode of op entrypoint](https://github.com/opctl/engine/issues/9)

## 0.1.1 - 2016-06-22
### Changed

- refactored to use [opspec sdk](https://github.com/opspec-io/sdk-golang)

### Fixed
- kill op run use case killing all ops
- [cannot run multiple ops with same name simultaneously](https://github.com/opctl/engine/issues/8)

### Removed

- add sub-op use case

## 0.1.0 - 2016-06-16
### Removed

- set op description use case
- add op use case
- list ops use case
