# Change Log

All notable changes to the spec will be documented in this file in
accordance with [keepachangelog.com](http://keepachangelog.com/)

## \[Unreleased]

### Added

- typed params; `dir`, `file`, `number`, `socket`, `string`
- `string` and `number` parameter constraints
- support for container calls
- `filter` to engine API `/event-stream` resource
- [support for private images](https://github.com/opspec-io/spec/issues/71)

### Changed

- op call changed from `string` to `object` w/ `ref`, `inputs`, and
  `outputs` attributes. To migrate, replace string value with object
  having `ref` attribute equal to existing string and add
  `inputs`/`outputs` values as applicable.
- String parameters must now be declared as an object:
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

- `serial`, `op`, and `parallel` calls
- nested calls (applicable to `serial` & `parallel` calls)
- json schema

### Changed

- params no longer support `type` attribute;
- `subOps` call; use new `op` call

## \[0.1.1] - 2016-08-03

## \[0.1.0] - 2016-07-18

