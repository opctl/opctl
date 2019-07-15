# Change Log

All notable changes to this project will be documented in this file in
accordance with
[![keepachangelog 1.0.0](https://img.shields.io/badge/keepachangelog-1.0.0-brightgreen.svg)](http://keepachangelog.com/en/1.0.0/)

## \[Unreleased]

### Added

- [Allow dynamically setting env vars of a container](https://github.com/opctl/specs/issues/250)
- [NotExists predicate](https://github.com/opctl/specs/issues/245)
- [Exists predicate](https://github.com/opctl/specs/issues/241)
- up to 10x disk performance improvement on OSX
- [Ability to specify custom node data dir](https://github.com/opctl/opctl/issues/449)
- [Allow Numbers for Container Ports](https://github.com/opctl/specs/issues/233)
- [Interpolate Container Name](https://github.com/opctl/specs/issues/232)
- [Conditional running](https://github.com/opctl/specs/issues/223)
- serialLoop call
- parallelLoop call

### Fixed

- [opctl ls on windows does not list anything](https://github.com/opctl/opctl/issues/460)
- [object & array initializers don't support multiline values](https://github.com/opctl/sdk-golang/issues/416)
- [errors from parallel calls not logged](https://github.com/opctl/opctl/issues/421)

### Removed

- [`stdOut` & `stdErr` attributes from container call.](https://github.com/opctl/specs/issues/231). Use files.
- `pkg` attribute in opcall; `ref` & `pullCreds` raised up a level, nesting within `pkg` unnecessary.

### Changed

- website/docs moved to [https://github.com/opctl/website](https://github.com/opctl/website)

## 0.1.24 - 2018-04-06

### Added

- [`opspec` property in op schema](https://github.com/opctl/specs/issues/20)
- Client back pressure in `GET event-stream` endpoint via `ack` query param
- [Support declaring SVG icon for op](https://github.com/opctl/specs/issues/139)
- [Support CommonMark for op & param descriptions](https://github.com/opctl/specs/issues/174)
- [Boolean type](https://github.com/opctl/specs/issues/195)
- [Support type, description, writeOnly, & title keywords in constraints of object params](https://github.com/opctl/specs/issues/196)
- [Support paths in object refs](https://github.com/opctl/specs/issues/170)
- Object & Array initializers
- Support referencing object properties via `object[propertyName]`
- Support referencing array items via `array[index]` or `array[-index]` to index from end of array
- [Interpolation escape syntax](https://github.com/opctl/specs/issues/193) by prefixing w/ a single backslash; i.e. `\$(not-a-ref)`
- [Unified data API](https://github.com/opctl/specs/issues/204)

### Deprecated

- `pkg` attribute in
  [op.yml.schema.json#/definitions/opCall](spec/op.yml.schema.json#/definitions/opCall); `ref` & `pullCreds` raised up a level, nesting within `pkg` unnecessary.
- `pkg` changed to `op` in [node-api.spec.yml#/components](spec/node-api.spec.yml#/components)
- [Deprecate pkgs API](https://github.com/opctl/specs/issues/205)
- `stdOut` & `stdErr` attributes from container call. Use files.

### Removed

- [References in/as expressions w/out explicit opener $( and closer )](https://github.com/opctl/specs/issues/184)

## 0.1.23 - 2018-01-15

### Added

- [opspec 0.1.6) Support declaring SVG icon for pkg](https://github.com/opctl/spec/issues/139)
- [opspec 0.1.6) Support CommonMark for pkg & param descriptions](https://github.com/opctl/spec/issues/174)

### Fixed

- coercion doesn't occur when de referencing scope object paths
- scope file path refs don't de reference

## 0.1.22 - 2017-11-05

### Added

- [Always pull container images when running ops](https://github.com/opctl/opctl/issues/362)

### Fixed

- [Auto node creation requires opctl bin within path](https://github.com/opctl/opctl/issues/363)

## 0.1.21 - 2017-10-01

### Added

- [Validation of outputs](https://github.com/opctl/opctl/issues/186)
- Remote pkg run
- Remote pkg validate
- [Type coercion](https://github.com/opctl/specs/issues/165)
- [Add /pkgs/{ref}/contents endpoints to node API](https://github.com/opctl/specs/issues/132)
- [Support binding strings &/or numbers to/from container files](https://github.com/opctl/specs/issues/131)
- [Add support for object param type](https://github.com/opctl/specs/issues/65)
- [Add support for array param type](https://github.com/opctl/specs/issues/160)

### Deprecated

- op.yml without `opspec` property
- References in/as expressions w/out explicit opener `$(` and closer `)`

### Fixed

- [Color flags not working](https://github.com/opctl/opctl/issues/278)
- [Race condition for non-cached pkgs](https://github.com/opctl/opctl/issues/253)
- [Binding pkg file/dir to sub op input doesn't work](https://github.com/opctl/opctl/issues/249)
- [Outputs not defaulted](https://github.com/opctl/opctl/issues/314)

### Removed

- `ref` attribute in
  [op.yml.schema.json#/definitions/opCall](spec/op.yml.schema.json#/definitions/opCall).
  Use new `pkg` attribute.
- `pullIdentity` & `pullSecret` attributes in
  [op.yml.schema.json#/definitions/containerCall](spec/op.yml.schema.json#/definitions/containerCall).
  Use new `pullCreds` attribute.

### Changed

- [node-api.spec.yml](spec/node-api.spec.yml) updated to OAS 3.0.0
  syntax

## 0.1.20 - 2017-06-23

### Fixed

- [Unexpected git capabilities encountered during pkg pull not gracefully handled](https://github.com/opctl/opctl/issues/257)

## 0.1.19 - 2017-06-05

### Added

- [Support using dir/file embedded in op as input/output param default](https://github.com/opctl/specs/issues/127)
- [Allow path expansion w/in sub op call inputs](https://github.com/opctl/specs/issues/120)
- [Allow string/number literals as sub op call inputs](https://github.com/opctl/specs/issues/121)
- [Implicitly bind env vars to in scope refs if names are identical](https://github.com/opctl/specs/issues/117)
- `pkg install` command
- [Validate file/dir inputs are valid files/dirs (respectively)](https://github.com/opctl/opctl/issues/175)
- [Fail fast during parallel call](https://github.com/opctl/opctl/issues/154)
- [Support since in event filter](https://github.com/opctl/opctl/issues/187)
- [Add github style heading links to web docs](https://github.com/opctl/opctl/issues/194)
- [Add copy code to clipboard link to web docs](https://github.com/opctl/opctl/issues/195)

### Deprecated

- `ref` attribute in
  [op.yml.schema.json#/definitions/opCall](spec/op.yml.schema.json#/definitions/opCall).
  Use new `pkg` attribute.
- `pullIdentity` & `pullSecret` attributes in
  [op.yml.schema.json#/definitions/containerCall](spec/op.yml.schema.json#/definitions/containerCall).
  Use new `pullCreds` attribute.

### Removed

- `pkg set` command

### Fixed

- [Killing a run (ctrl+c) from powershell hangs](https://github.com/opctl/opctl/issues/199)
- [Network creation race condition](https://github.com/opctl/opctl/issues/190)
- [Param defaults w/ values equal to type default are not defaulted](https://github.com/opctl/opctl/issues/185)
- [stdOut/stdErr output race condition](https://github.com/opctl/opctl/issues/174)
- [Unable to run ops w/ containers if using docker 4 windows](https://github.com/opctl/opctl/issues/200)


## 0.1.18 - 2017-03-28

### Changed

- [Don't recreate node on self-update](https://github.com/opctl/opctl/issues/169)

### Fixed

- [Multiple opctl networks created leading to lack of inter-container connectivity](https://github.com/opctl/opctl/issues/167)

## 0.1.16 - 2017-03-26

### Fixed

- [Outputs internal to op call graph not initialized](https://github.com/opctl/opctl/issues/165)

## 0.1.15 - 2017-03-23

### Added

- Add `node` command w/ `create` and `kill` subcommands
- [Add ability to override default (`.opspec`) package location for `pkg set`, `pkg create`, `run`, and `ls` commands](https://github.com/opctl/opctl/issues/44)
- [Add output coloring](https://github.com/opctl/opctl/issues/49)
- Add input validation
- Added package validation via `pkg validate` command & before `run`
- Add `pkg` command w/ `validate`, `set`, `create` subcommands
- typed params; `dir`, `file`, `number`, `socket`, `string`
- `string` and `number` parameter constraints
- support for container calls
- `filter` to node API `/events/stream` resource
- [support for private images](https://github.com/opctl/specs/issues/71)

### Changed

- op call changed from `string` to `object` w/ `ref`, `inputs`, and
  `outputs` attributes. To migrate, replace string value with object
  having `ref` attribute equal to existing string and pass
  `inputs`/`outputs` as applicable.
- String parameters must now be declared as an object:

```yaml
paramName:
      string:
        description: ...
        # and so on...
```

### Removed

- `docker-compose.yml`; replaced with container calls
- collections
- bubbling of default collection lookup
- support for < [opspec 0.1.3](https://opspec.io)
- `collection` command

## 0.1.10 - 2016-11-21

### Added

- [Add support for "default" input values](https://github.com/opctl/opctl/issues/41)

## 0.1.9 - 2016-11-06

### Added

- `serial`, `op`, and `parallel` calls
- nested calls (applicable to `serial` & `parallel` calls)
- json schema

### Changed

- refactored to use [sdks/go](https://github.com/opctl/opctl/sdks/go)
- params no longer support `type` attribute;
- `subOps` call; use new `op` call

### Fixed

- [Emitted ContainerStd*WrittenToEvent.Data Incomplete](https://github.com/opctl/opctl/issues/32)

## 0.1.8 - 2016-09-09

### Added

- support for [opspec 0.1.2](https://opspec.io)

### Fixed

- [failure of serial operation run does not immediately fail all following operations](https://github.com/opctl/cli/issues/5)

### Removed

- support for < [opspec 0.1.2](https://opspec.io)

## 0.1.7 - 2016-09-02

### Fixed

- [opctl does not wait for parallel op containers to die before returning](https://github.com/opctl/cli/issues/8)
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
