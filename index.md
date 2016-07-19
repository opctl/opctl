# Table of Contents

- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Terminology](#terminology)
- [Op](#op)
    - [Op Bundle](#op-bundle)
- [Collection](#collection)
    - [Collection Bundle](#collection-bundle)
    - [Default Collection](#default-collection)
- [Registry](#registry)
    - [Registry API](#registry-api)
- [Engine](#engine)
    - [Engine API](#engine-api)

# Introduction

## Purpose

Opspec is a framework for specifying, distributing, and executing ops
(operations).

Primary concerns of opspec are to make operations:

- composable
- containerized
- distributable
- fully specified
- versionable

## Terminology

### FILE_SAFE_NAME

a string matching the regex `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`

### FILE_BUNDLE

a file directory with a defined structure allowing related files to be
grouped together as a conceptually single item (see <a
href="https://en.wikipedia.org/wiki/Bundle_(OS_X)">Bundle_(OS_X)</a> for
similar usage)

### MUST/MAY

as defined by [RFC 2119](https://tools.ietf.org/html/rfc2119)

# Op

An op is a process or task.


## Op Bundle

Ops are defined in the form of a [FILE_BUNDLE](#file_bundle) adhering to
the following criteria:

- [MUST](#mustmay) contain an [op file](op-file.md) at the root of the
  op directory.
- [MUST](#mustmay) contain an
  [op docker compose file](op-docker-compose-file.md) at the root of the
  op directory.

**example op bundle**  

```TEXT
  |-- op.yml
  |-- docker-compose.yml
  ... (op specific files/dirs)
```

Once defined, **op bundles**, can be distributed via any traditional
means of file transfer.

# Collection

A collection groups one or more [op](#op)s together


## Collection Bundle

Collections are defined in the form of a [FILE_BUNDLE](#file_bundle)
adhering to the following criteria:

- [MUST](#mustmay) contain an [collection file](collection-file.md) at
  the root of the collection directory.
- [MAY](#mustmay) contain one or more [op](#op)s at the root of the
  collection directory.

**example collection bundle**  

```TEXT
  |-- collection.yml
  ... (embedded op bundles)
```

Once defined, **collection bundles**, can be published to

## Default Collection

A default collection may be designated for a directory according to the
following criteria:

- [MUST](#mustmay) be contained by an `.opspec` directory

In the event a default collection is not present in a directory, its
nearest ancestor [MUST](#mustmay) be used as the effective default


# Registry

A registry allows publishing [op bundle](#op-bundle)s &
[collection bundle](#collection-bundle)s to **repos**, enabling them to
be centrally discovery & consumed.

## Registry API

Registries [MUST](#mustmay) implement the
[registry api](registry-oai_spec.yaml)


# Engine

An engine runs ops.

## Engine API

Engines [MUST](#mustmay) implement the
[engine api](engine-oai_spec.yaml)
