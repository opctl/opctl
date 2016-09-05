# Change Log

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added 

- typed params; `string`
- typed op run statements; `serial`, `op`, and `parallel`
- nested op run statement support (applies to `serial` & `parallel` op
  run statements)
- json schema

### Changed

- params no longer support `type` attribute; use new `string` param type
- op run statement no longer supports `subOps` attribute; use new `op`
  run statement type

## [0.1.1] - 2016-08-03

## [0.1.0] - 2016-07-18

