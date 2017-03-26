# Change Log

All notable changes to the spec will be documented in this file in
accordance with [keepachangelog.com](http://keepachangelog.com/)

## \[Unreleased]

### Deprecated

- `ref` attribute in
  [package-manifest.schema.json#/definitions/opCall](spec/package-manifest.schema.json#/definitions/opCall).
  Use new `pkg` attribute.

### Removed

- `pullIdentity` & `pullSecret` attributes in
  [package-manifest.schema.json#/definitions/containerCall](spec/package-manifest.schema.json#/definitions/containerCall).
  Use new `pullAuth` attribute.

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

