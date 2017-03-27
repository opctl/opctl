# Change Log

All notable changes to this project will be documented in this file in
accordance with [keepachangelog.com](http://keepachangelog.com/)

## 0.1.16 - 2017-03-26

### Fixed

- [Outputs internal to op call graph not initialized](https://github.com/opctl/opctl/issues/165)

## 0.1.15 - 2017-03-23

### Added

- Add `node` command w/ `create` and `kill` subcommands
- [Add ability to override default (`.opspec`) package location for `pkg set`, `pkg create`, `run`, and `ls` commands](https://github.com/opctl/opctl/issues/44)
- [Add opspec 0.1.3 support](https://github.com/opctl/opctl/issues/48)
- [Add output coloring](https://github.com/opctl/opctl/issues/49)
- Add input/output constraint validation
- Added package validation via `pkg validate` command & before `run`
- Add `pkg` command w/ `validate`, `set`, `create` subcommands

### Removed

- support for < [opspec 0.1.3](https://opspec.io)
- `collection` command
- `op` command

## 0.1.10 - 2016-11-21

### Added

- [Add support for "default" input values](https://github.com/opctl/opctl/issues/41)

## 0.1.9 - 2016-11-06

### Changed

- refactored to use
  [opspec sdk](https://github.com/opspec-io/sdk-golang)

### Fixed

- [Emitted ContainerStd*WrittenToEvent.Data Incomplete](https://github.com/opctl/opctl/issues/32)

## 0.1.8 - 2016-09-09

### Added

- support for [opspec 0.1.2](https://opspec.io)

### Fixed

- [failure of serial operation run does not immediately fail all following operations](https://github.com/opspec-io/cli/issues/5)

### Removed

- support for < [opspec 0.1.2](https://opspec.io)

## 0.1.7 - 2016-09-02

### Fixed

- [opctl does not wait for parallel op containers to die before returning](https://github.com/opspec-io/cli/issues/8)
- [Many parallel ops crash engine](https://github.com/opctl/opctl/issues/17)

## 0.1.6 - 2016-08-21

### Fixed

- OpEnded event not sent on `Failed` outcome

## 0.1.5 - 2016-08-02

### Added

- support for [opspec 0.1.1](https://opspec.io)

### Removed

- support for [opspec 0.1.0](https://opspec.io)

## 0.1.4 - 2016-07-20

### Added

- normalization of windows paths if provided to op run

## 0.1.3 - 2016-07-09

### Added

- [Support new opspec subop `isParallel` flag](https://github.com/opctl/opctl/issues/11)

### Fixed

- [Unable to simultaneously run multiple ops from same collection](https://github.com/opctl/opctl/issues/10)

## 0.1.2 - 2016-06-22

### Fixed

- [Missleading `variable is not set` message on op finish](https://github.com/opctl/opctl/issues/5)
- [Engine not observing exitcode of op entrypoint](https://github.com/opctl/opctl/issues/9)

## 0.1.1 - 2016-06-22

### Changed

- refactored to use opspec engine sdk

### Fixed

- kill op run use case killing all ops
- [cannot run multiple ops with same name simultaneously](https://github.com/opctl/opctl/issues/8)

### Removed

- add sub-op use case

## 0.1.0 - 2016-06-16

### Removed

- set op description use case
- add op use case
- list ops use case
