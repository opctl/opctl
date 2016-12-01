# Change Log

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- `dir`, `file`, and `netSocket` parameter types
- `pattern`, `minLength`, and `maxLength` validation attributes for
  `string` parameters
- `container` run declarations
- `designated collections`

### Changed

- Op run declaration changed from string to an object with `ref`, `arg`,
  and `result` attributes. To migrate, replace string value with object
  having `ref` attribute equal to existing string.
- String inputs must now be declared as a key value pair where the key
  is `string` and the value is its attributes. To migrate, just turn
  existing inputs into a key value pair where the key is `string` and
  existing attributes are the value.

### Removed
- `default collections`. Replaced with `designated collections`

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

