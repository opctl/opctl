# Change Log

All notable changes will be documented in this file in accordance with
[![keepachangelog 1.0.0](https://img.shields.io/badge/keepachangelog-1.0.0-brightgreen.svg)](http://keepachangelog.com/en/1.0.0/)

## \[Unreleased]

## \[0.1.5] - 2017-09-27

### Added

- [Type coercion](https://github.com/opspec-io/spec/issues/165)
- [Add /pkgs/{ref}/contents endpoints to node API](https://github.com/opspec-io/spec/issues/132)
- [Support binding strings &/or numbers to/from container files](https://github.com/opspec-io/spec/issues/131)
- [Add support for object param type](https://github.com/opspec-io/spec/issues/65)
- [Add support for array param type](https://github.com/opspec-io/spec/issues/160)

### Deprecated

- pkg fs & scope references in/as op call args without `$(ref)` reference syntax

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

- [Support using pkg dir/file as input/output param default](https://github.com/opspec-io/spec/issues/127)
- [Allow path expansion w/in sub op call inputs](https://github.com/opspec-io/spec/issues/120)
- [Allow string/number literals as sub op call inputs](https://github.com/opspec-io/spec/issues/121)
- [Implicitly bind env vars to in scope refs if names are identical](https://github.com/opspec-io/spec/issues/117)

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
- [support for private images](https://github.com/opspec-io/spec/issues/71)

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
- collections; replaced w/ package resolution
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

