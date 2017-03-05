# Change Log

All notable changes to the spec will be documented in this file.

## \[Unreleased]

### Added

- `dir`, `file`, `number`, `socket` parameter types
- `string` and `number` parameter constraints
- support for container calls
- `filter` to engine API `/event-stream` resource
- [support for private images](https://github.com/opspec-io/spec/issues/71)

### Changed

- op call changed from string to object with `ref`, `inputs`, and `outputs`
  attributes. To migrate, replace string value with object having `ref`
  attribute equal to existing string and add `inputs`/`outputs` values as
  applicable.
- String parameters must now be declared as an object in form:
  ```yaml
  paramName:
    string:
      description: ...
      # and so on... 
  ```

### Removed

- bubbling of default collection lookup

## \[0.1.2] - 2016-09-10

### Added 

- typed run declarations; `serial`, `op`, and `parallel`
- nested run declaration support (applies to `serial` & `parallel` run
  declarations)
- json schema

### Changed

- params no longer support `type` attribute;
- run declaration no longer supports `subOps` attribute; use new `op`
  run declaration type

## \[0.1.1] - 2016-08-03

## \[0.1.0] - 2016-07-18

