# Change Log

All notable changes will be documented in this file in accordance with
[![keepachangelog 1.0.0](https://img.shields.io/badge/keepachangelog-1.0.0-brightgreen.svg)](http://keepachangelog.com/en/1.0.0/)

## [Unreleased]

### Added

- [NotExists predicate](https://github.com/opctl/specs/issues/245)
- [Exists predicate](https://github.com/opctl/specs/issues/241)
- [Allow Numbers & Implicit Binding On Container Ports](https://github.com/opctl/specs/issues/233)
- [Interpolate Container Name](https://github.com/opctl/specs/issues/232)
- [Conditional running](https://github.com/opctl/specs/issues/223)
- [Looping](https://github.com/opctl/specs/issues/207)

### Removed

- [`stdOut` & `stdErr` attributes from container call.](https://github.com/opctl/specs/issues/231). Use files.
- `pkg` attribute in
  [op.yml.schema.json#/definitions/opCall](spec/op.yml.schema.json#/definitions/opCall); `ref` & `pullCreds` raised up a level, nesting within `pkg` unnecessary.

## \[0.1.6] - 2018-04-05

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

## \[0.1.5] - 2017-09-27

### Added

- [Type coercion](https://github.com/opctl/specs/issues/165)
- [Add /pkgs/{ref}/contents endpoints to node API](https://github.com/opctl/specs/issues/132)
- [Support binding strings &/or numbers to/from container files](https://github.com/opctl/specs/issues/131)
- [Add support for object param type](https://github.com/opctl/specs/issues/65)
- [Add support for array param type](https://github.com/opctl/specs/issues/160)

### Deprecated

- op.yml without `opspec` property
- References in/as expressions w/out explicit opener `$(` and closer `)`

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


## \[0.1.4] - 2017-06-04

### Added

- [Support using dir/file embedded in op as input/output param default](https://github.com/opctl/specs/issues/127)
- [Allow path expansion w/in sub op call inputs](https://github.com/opctl/specs/issues/120)
- [Allow string/number literals as sub op call inputs](https://github.com/opctl/specs/issues/121)
- [Implicitly bind env vars to in scope refs if names are identical](https://github.com/opctl/specs/issues/117)

### Deprecated

- `ref` attribute in
  [op.yml.schema.json#/definitions/opCall](spec/op.yml.schema.json#/definitions/opCall).
  Use new `pkg` attribute.
- `pullIdentity` & `pullSecret` attributes in
  [op.yml.schema.json#/definitions/containerCall](spec/op.yml.schema.json#/definitions/containerCall).
  Use new `pullCreds` attribute.

## \[0.1.3] - 2017-03-06

### Added

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

## \[0.1.2] - 2016-09-10

### Added

- `serial`, `op`, and `parallel` calls
- nested calls (applicable to `serial` & `parallel` calls)
- json schema

### Changed

- params no longer support `type` attribute;
- `subOps` call; use new `op` call

## \[0.1.1] - 2016-08-03

## \[0.1.0] - 2016-07-18

