# Change Log

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- `dir`, `file`, and `netSocket` parameter types
- `constraints` attribute for `string` and `netSocket` parameter types
- [./schema/call-graph.json#definitions/containerCall](schema/call-graph.json#definitions/containerCall)

### Changed

- Rename from `run` to `call graph` for consistency with established
  terminology.
- [./schema/call-graph.json#definitions/opCall](schema/call-graph.json#definitions/opCall)
  changed from string to object with `ref`, `arg`, and `result`
  attributes. To migrate, replace string value with object having `ref`
  attribute equal to existing string and add `arg`/`result` values as
  applicable.
- String inputs must now be declared as a key value pair where the key
  is `string` and the value is its attributes. To migrate, just turn
  existing inputs into a key value pair where the key is `string` and
  existing attributes are the value.

### Removed

- bubbling of default collection lookup

## [0.1.2] - 2016-09-10

### Added 

- typed run declarations; `serial`, `op`, and `parallel`
- nested run declaration support (applies to `serial` & `parallel` run
  declarations)
- json schema

### Changed

- params no longer support `type` attribute;
- run declaration no longer supports `subOps` attribute; use new `op`
  run declaration type

## [0.1.1] - 2016-08-03

## [0.1.0] - 2016-07-18

