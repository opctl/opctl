# Changelog

All notable changes to this project will be documented in this file in
accordance with
[![keepachangelog 1.0.0](https://img.shields.io/badge/keepachangelog-1.0.0-brightgreen.svg)](http://keepachangelog.com/en/1.0.0/)

## [0.1.72] - 2025-04-13

### Changed

- CLI command examples now are useable when possible and comments now use full sentences

### Fixed

- CLI command help no longer includes empty `Examples` section when no examples exist
- CLI `op validate` command broken for remote ops

## [0.1.71] - 2025-04-10

### Added

- Examples in CLI help/usage
- Autogenerate CLI docs for website from actual CLI

### Changed

- CLI version command now outputs to stdout (instead of stderr)

### Fixed

- `sudo: opctl command not found` error if `opctl` not on PATH

## [0.1.70] - 2025-04-02

### Added

- Support default auth for pulling ops via `opctl auth add [...]`

## [0.1.69] - 2025-04-02

### Added

- Support for specifying image.platform.arch in container calls

## [0.1.68] - 2025-03-22

### Fixed

- Fix accessing containers by name from opctl nodes stops working due to prior registrations not being cleaned up on container exit

## [0.1.67] - 2025-03-08

### Fixed

- Fix accessing Docker For Mac 4.39.0+ containers by name from opctl nodes

## [0.1.66] - 2025-02-26

### Added

- Support for environment variables for each CLI command options and arguments

## [0.1.65] - 2025-02-10

### Added

- Automatic handling of mDNSResponder port conflict on OSX

### Fixed

- Binding opctl DNS to a non-standard port (e.g. 54)

## [0.1.64] - 2025-02-10

### Added

- Leverage native privilege escalation handling for self-update

## [0.1.63] - 2025-02-08

### Added

- Native privilege escalation handling (i.e. no more calling with sudo)

## [0.1.62] - 2025-02-05

### Fixed

- Put in a short term fix so that when `pullCreds` is used on Linux against
something other than `docker.io`, things don't error out any more.

## [0.1.60] - 2024-11-24

### Fixed

- Fix a regression in the 'opctl op install' command introduced in 0.1.48 which caused it to be much slower than previous
- Fix a potential race condition encountered when pulling/using multiple remote ops from the same repo at once

## [0.1.59] - 2024-11-17

### Fixed

- `No such container: opctl_[...]` errors if container exits too fast (race condition)

## [0.1.58] - 2024-11-13

### Fixed

- Fix executable bit not maintained within remote ops

## [0.1.57] - 2024-11-12

### Changed

- Don't log "not found in graph" to logs (noisy without value)

### Fixed

- Fix 'opctl node kill' panics if you've never run an op
- Rate limit errors surfacing as "image not found" errors
- Fixed file initialization results in incorrect (empty) file data

## [0.1.56] - 2024-11-11

### Added

- Access containers by their name from opctl nodes

### Deprecated

- `container.ports`; access containers by name instead

## [0.1.56-alpha.1] - 2024-11-08

### Added

- Access containers by their name from opctl nodes

### Deprecated

- `container.ports`; access containers by name instead

## [0.1.56-alpha.0] - 2024-10-30

### Added

- Access containers by their name from opctl nodes

### Deprecated

- `container.ports`; access containers by name instead

## [0.1.55] - 2024-09-04

### Fixed

- Fix long time to resolve non-existent op

## [0.1.53] - 2024-04-18

### Fixed

- A few dependency updates
  - in the Golang SDK the only breaking upgrade is Docker v23 -> v25.0.3+incompatible
  - in the JS SDK nock (a devDependency) was updated from 9 -> 13.
  - in the JS SDK react-ace was updated from 9 -> 11

### Removed

- Deleted the React SDK because it is unused

## [0.1.52] - 2023-06-06

### Fixed

- [Escaped references are not escaped if an unescaped '$' exists prior in the string](https://github.com/opctl/opctl/issues/1063)

## [0.1.51] - 2023-06-01

### Added

- [Variable reference as `container.cmd`](https://github.com/opctl/opctl/issues/1064)

### Fixed

- Simultaneously defining a default plus constraints on an object or array input results in a validation error

## [0.1.50] - 2023-05-01

### Added

- Automatic GPU detection/passthrough for docker runtime

### Fixed

- `opctl ls` renders ops with no description as blank line
- `lte` predicate
- In some cases, stacktraces logged as bytes rather than strings
- Containers aren't automatically removed after exit for docker runtime
- In some cases, `Runtime error: invalid memory address or nil pointer dereference` error occurs
- CLI not logging errors from the root op

## [0.1.49] - 2022-09-21

### Added

- opspec now supports gt, gte, lt, lte predicates
- `opctl node kill` will now stop and remove any opctl managed containers
- introduced `opctl node delete` command which "Deletes a node. This is destructive! all node data including auth, caches, and operation state will be permanently removed."

### Changed

- upgrading to this version from prior versions is destructive! all node data including auth, caches, and operation state will be permanently removed.
- K8s container runtime now explicitly deletes terminated pods

### Fixed

- [Node locking mechanism doesn't ensure process is opctl](https://github.com/opctl/opctl/issues/913)
- [CLI no longer logs errors occurring in parallel calls](https://github.com/opctl/opctl/issues/1032)

## [0.1.48] - 2021-08-13

### Added

- When running an op via opctl run, display progress via a live call graph
- When running an op via opctl run, prefix log lines emitted by workloads with their op id & ref
- Basic support for sending local files and directories to remote nodes when using the API client
- [Allow defining description on call graph nodes](https://github.com/opctl/opctl/issues/900)

### Changed

- Self-update now uses github releases instead of equinox.io
- API now limits request body to 40Mb
- [Improved error output when op resolution fails. You'll now see a list of resolutions tried and why each failed.](https://github.com/opctl/opctl/pull/883)
- [More consistent error messaging formats](https://github.com/opctl/opctl/pull/885)
- [Detect invalid op output names](https://github.com/opctl/opctl/issues/798)
- [Allow using type initializers in input/output defaults](https://github.com/opctl/opctl/issues/957)
- [Deprecated absolute paths as file/dir input/output defaults](https://github.com/opctl/opctl/issues/957)
- [Deprecated op output binding syntax; use same syntax as binding inputs](https://github.com/opctl/opctl/issues/721)
- [Deprecated param.<datatype>.description; use param.description](https://github.com/opctl/opctl/issues/898)
- [Docker images will only be pulled if using the `latest` tag (or untagged) or have not been pulled previously](https://github.com/opctl/opctl/issues/920)
- Go SDK models now use DataRef rather than PkgRef

### Fixed

- [vscode intellisense error](https://github.com/opctl/opctl/issues/615)

### Removed

- pkgs API endpoint; use data API endpoint
- Windows build; use linux build via WSL 2 instead

## [0.1.47] - 2021-01-22

### Added

- [Improve CLI prompts for username and password](https://github.com/opctl/opctl/issues/745)

### Fixed

- Dir initializer doesn't initialize more than one child entry

## [0.1.46] - 2021-01-04

### Fixed

- [container.envVars string double interpreted](https://github.com/opctl/opctl/issues/849)

## [0.1.45] - 2020-11-17

### Fixed

- Calls killed by needs declaration exiting non-zero

## [0.1.44] - 2020-11-16

### Added

- [Dir Initializer Syntax](https://github.com/opctl/opctl/issues/500)

### Changed

- [Opspec) use relative paths for file/dir refs](https://github.com/opctl/opctl/issues/834)
- [Make input/output binding when calling ops consistent](https://github.com/opctl/opctl/issues/721)

### Fixed

- Certain child call errors not shown.

## [0.1.43] - 2020-11-04

### Fixed

- ParallelLoop loop iteration vars sometimes get set to values from other iterations.

## [0.1.42] - 2020-11-03

### Added

- [Add ability to add auth to opctl for OCI image registries](https://github.com/opctl/opctl/issues/823)
- [Better messages for parallel/parallelLoop child errors](https://github.com/opctl/opctl/issues/827)

### Changed

- [Make OpKill Event Driven](https://github.com/opctl/opctl/issues/809)
- [Remove CallKilled event](https://github.com/opctl/opctl/issues/810)
- [Remove ContainerExited event](https://github.com/opctl/opctl/issues/825)
- [Remove OpErred event](https://github.com/opctl/opctl/issues/812)
- [Remove Event suffix from event types](https://github.com/opctl/opctl/issues/814)
- [Rename types from SCG/DCG](https://github.com/opctl/opctl/issues/816)

### Fixed

- [Gracefully handle docker restarts](https://github.com/opctl/opctl/issues/678)
- [Running an op should never kill a node](https://github.com/opctl/opctl/issues/756)

## [0.1.41] - 2020-06-03

### Changed

- [Listen on localhost by default](https://github.com/opctl/opctl/issues/738)

## [0.1.39] - 2020-05-04

### Fixed

- ["manifest has unsupported version: 4" errors on newer versions of opctl](https://github.com/opctl/opctl/issues/768)

## [0.1.38] - 2020-05-03

### Changed

- [Make opctl ls error if invalid ops are encountered](https://github.com/opctl/opctl/issues/708)
- [Return ref instead of name from opctl ls](https://github.com/opctl/opctl/issues/634)

### Fixed

- [Nested ops can't be referenced using relative path](https://github.com/opctl/opctl/issues/762)
- [Inconsistent behavior when running locally installed vs remotely referenced ops.](https://github.com/opctl/opctl/issues/732)

## [0.1.37] - 2020-05-03

### Added

- [ui subcommand to open webui](https://github.com/opctl/opctl/issues/758)
- Render op icons in UI
- Automatically expand mount ancestors in explorer UI
- Make call bounding box extend from call summary rather than start below in UI
- Remove extraneous lines extending from top and bottom of parallel call in UI

## [0.1.35] - 2020-04-29

### Changed

- [Stop logging "Replaying from value pointer: {Fid:0 Len:0 Offset:0}"](https://github.com/opctl/opctl/issues/754)

## [0.1.34] - 2020-04-23

### Added

- [UI: visualize referenced ops](https://github.com/opctl/opctl/issues/739)

## [0.1.33] - 2020-04-20

### Fixed

- [Nonexistent sub dirs bound to containers aren't sync'd](https://github.com/opctl/opctl/issues/725)
- [image.ref with multi-variable templated string not working since v0.1.28](https://github.com/opctl/opctl/issues/722)

## [0.1.32] - 2020-04-16

### Added

- [Prefix opctl managed container names with opctl\_](https://github.com/opctl/opctl/issues/735)

### Fixed

- variable reference validation triggered for valid refs

## [0.1.31] - 2020-04-15

### Fixed

- variable reference validation triggered for valid refs
- failure interpreting needed call panics

## [0.1.30] - 2020-04-15

### Added

- [Named calls and needs](https://github.com/opctl/opctl/issues/643)

## [0.1.29] - 2020-04-02

### Fixed

- [Running `op install` twice wipes out op file contents](https://github.com/opctl/opctl/issues/718)

## [0.1.28] - 2020-03-26

### Added

- UI: Workspace page (explorer, op visualizer with pan/zoom)
- [Support in scope dir as op](https://github.com/opctl/opctl/issues/646)
- Liveness method to node API Client
- Variable reference as `image.ref`.

### Changed

- When daemonizing opctl node, parent process env vars no longer inherited by daemonized process. This for example thwarts Jenkins ProcessTreeKiller's killing abilities.

### Deprecated

- `image.src`; use `image.ref`

### Fixed

- API Liveness endpoint incorrectly returning 404

### Removed

- UI: Events/Op/Vars pages

## [0.1.27] - 2020-02-04

### Fixed

- Object initializers passed as inputs to constrained parameters don't pass validation

## [0.1.26] - 2020-01-30

### Added

- [Support in scope dir as container image](https://github.com/opctl/opctl/issues/498)
- [Pass thru errors encountered when cli auto daemonizes a node](https://github.com/opctl/opctl/issues/368)
- [Allow Interpolating Container `workDir`](https://github.com/opctl/opctl/issues/648)
- `Container.sockets` bindings with variable reference syntax i.e. `/my/socket: $(mySocket)`

### Changed

- Don't cleanup OPCTL_DATA_DIR on node creation.

### Deprecated

- `Container.sockets` bindings without variable reference syntax i.e. instead of `/my/socket: mySocket`, use `/my/socket: $(mySocket)`.

### Fixed

- [Referencing Non Directories As Directories Hangs](https://github.com/opctl/opctl/issues/637)
- Implicit inputs not coerced
- Results of Serial Call Children Running For > 10s Ignored
- [Child Op Call Inputs Not Required](https://github.com/opctl/opctl/pull/665)

### Removed

- Remove support for `.` in op parameter names (to avoid ambiguity between referencing object properties)

## [0.1.25] - 2019-07-13

### Added

- Allow dynamically setting env vars of a container
- NotExists predicate
- Exists predicate
- up to 10x disk performance improvement on OSX
- Ability to specify custom node data dir
- Allow Numbers for Container Ports
- Interpolate Container Name
- Conditional running
- serialLoop call
- parallelLoop call

### Fixed

- opctl ls on windows does not list anything
- object & array initializers don't support multiline values
- errors from parallel calls not logged

### Removed

- `stdOut` & `stdErr` attributes from container call. Use files.
- `pkg` attribute in opcall; `ref` & `pullCreds` raised up a level, nesting within `pkg` unnecessary.

### Changed

- website/docs moved to [https://github.com/opctl/website](https://github.com/opctl/website)

## [0.1.24] - 2018-04-06

### Added

- `opspec` property in op schema
- Client back pressure in `GET event-stream` endpoint via `ack` query param
- Support declaring SVG icon for op
- Support CommonMark for op & param descriptions
- Boolean type
- Support type, description, writeOnly, & title keywords in constraints of object params
- Support paths in object refs
- Object & Array initializers
- Support referencing object properties via `object[propertyName]`
- Support referencing array items via `array[index]` or `array[-index]` to index from end of array
- Interpolation escape syntax by prefixing w/ a single backslash; i.e. `\$(not-a-ref)`
- Unified data API

### Deprecated

- `pkg` attribute in opcall; `ref` & `pullCreds` raised up a level, nesting within `pkg` unnecessary
- Deprecate pkgs API
- `stdOut` & `stdErr` attributes from container call. Use files.

### Removed

- References in/as expressions w/out explicit opener $( and closer )

## [0.1.23] - 2018-01-15

### Added

- opspec 0.1.6) Support declaring SVG icon for pkg
- opspec 0.1.6) Support CommonMark for pkg & param descriptions

### Fixed

- coercion doesn't occur when de referencing scope object paths
- scope file path refs don't de reference

## [0.1.22] - 2017-11-05

### Added

- Always pull container images when running ops

### Fixed

- Auto node creation requires opctl bin within path

## [0.1.21] - 2017-10-01

### Added

- Validation of outputs
- Remote pkg run
- Remote pkg validate
- Type coercion
- Add /pkgs/{ref}/contents endpoints to node API
- Support binding strings &/or numbers to/from container files
- Add support for object param type
- Add support for array param type

### Deprecated

- op.yml without `opspec` property
- References in/as expressions w/out explicit opener `$(` and closer `)`

### Fixed

- [Color flags not working](https://github.com/opctl/opctl/issues/278)
- [Race condition for non-cached pkgs](https://github.com/opctl/opctl/issues/253)
- [Binding pkg file/dir to sub op input doesn't work](https://github.com/opctl/opctl/issues/249)
- [Outputs not defaulted](https://github.com/opctl/opctl/issues/314)

### Removed

- `ref` attribute in opcall
  Use new `pkg` attribute.
- `pullIdentity` & `pullSecret` attributes in container call.
  Use new `pullCreds` attribute.

### Changed

- api schema updated to OAS 3.0.0
  syntax

## [0.1.20] - 2017-06-23

### Fixed

- [Unexpected git capabilities encountered during pkg pull not gracefully handled](https://github.com/opctl/opctl/issues/257)

## [0.1.19] - 2017-06-05

### Added

- Support using dir/file embedded in op as input/output param default
- Allow path expansion w/in sub op call inputs
- Allow string/number literals as sub op call inputs
- Implicitly bind env vars to in scope refs if names are identical
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

## [0.1.18] - 2017-03-28

### Changed

- [Don't recreate node on self-update](https://github.com/opctl/opctl/issues/169)

### Fixed

- [Multiple opctl networks created leading to lack of inter-container connectivity](https://github.com/opctl/opctl/issues/167)

## [0.1.16] - 2017-03-26

### Fixed

- [Outputs internal to op call graph not initialized](https://github.com/opctl/opctl/issues/165)

## [0.1.15] - 2017-03-23

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
- support for private images

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

## [0.1.10] - 2016-11-21

### Added

- [Add support for "default" input values](https://github.com/opctl/opctl/issues/41)

## [0.1.9] - 2016-11-06

### Added

- `serial`, `op`, and `parallel` calls
- nested calls (applicable to `serial` & `parallel` calls)
- json schema

### Changed

- refactored to use [sdks/go](https://github.com/opctl/opctl/sdks/go)
- params no longer support `type` attribute;
- `subOps` call; use new `op` call

### Fixed

- [Emitted ContainerStd\*WrittenToEvent.Data Incomplete](https://github.com/opctl/opctl/issues/32)

## [0.1.8] - 2016-09-09

### Added

- support for [opspec 0.1.2](https://opspec.io)

### Fixed

- [failure of serial operation run does not immediately fail all following operations](https://github.com/opctl/cli/issues/5)

### Removed

- support for < [opspec 0.1.2](https://opspec.io)

## [0.1.7] - 2016-09-02

### Fixed

- [opctl does not wait for parallel op containers to die before returning](https://github.com/opctl/cli/issues/8)
- [Many parallel ops crash engine](https://github.com/opctl/opctl/issues/17)

## [0.1.6] - 2016-08-21

### Fixed

- OpEnded event not sent on `Failed` outcome

## [0.1.5] - 2016-08-02

### Added

- support for [opspec 0.1.1](https://opspec.io)

### Removed

- support for [opspec 0.1.0](https://opspec.io)

## [0.1.4] - 2016-07-20

### Added

- normalization of windows paths if provided to op run

## [0.1.3] - 2016-07-09

### Added

- [Support new opspec subop `isParallel` flag](https://github.com/opctl/opctl/issues/11)

### Fixed

- [Unable to simultaneously run multiple ops from same collection](https://github.com/opctl/opctl/issues/10)

## [0.1.2] - 2016-06-22

### Fixed

- [Missleading `variable is not set` message on op finish](https://github.com/opctl/opctl/issues/5)
- [Engine not observing exitcode of op entrypoint](https://github.com/opctl/opctl/issues/9)

## [0.1.1] - 2016-06-22

### Changed

- refactored to use opspec engine sdk

### Fixed

- kill op run use case killing all ops
- [cannot run multiple ops with same name simultaneously](https://github.com/opctl/opctl/issues/8)

### Removed

- add sub-op use case

## [0.1.0] - 2016-06-16

### Removed

- set op description use case
- add op use case
- list ops use case
