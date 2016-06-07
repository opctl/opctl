# Table of Contents
- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Terminology](#terminology)
- [Op](#op)
- [Op Collection](#op-collection)
- [Default Op Collection](#default-op-collection)

# Introduction

## Purpose
Op Spec is an op specification format.

Primary concerns of op spec are to make operations:
- portable across runtimes
- distributable
- discoverable
- versionable

## Terminology

### FILE_SAFE_NAME
a string matching the regex `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`

### FILE_BUNDLE
a file directory with a defined structure allowing related files to be grouped 
together as a conceptually single item 
(see <a href="https://en.wikipedia.org/wiki/Bundle_(OS_X)">Bundle_(OS_X)</a> for similar usage)

### MUST/MAY
as defined by [RFC 2119](https://tools.ietf.org/html/rfc2119)

# Op
An op is a process or task. 

Ops are defined in the form of a [FILE_BUNDLE](#file_bundle) adhering to the 
following criteria:

- [MUST](#mustmay) contain an [op file](op-file.md) at the root of the op directory.
- [MUST](#mustmay) contain an [op docker compose file](./op-docker-compose-file.md) 
at the root of the op directory.

**op file bundle**  
```TEXT
  |-- op.yml
  |-- docker-compose.yml
  ... (op specific files/dirs)
```

# Op Collection
An op collection groups one or more [op](#op)s together

Op collections are defined in the form of a [FILE_BUNDLE](#file_bundle) adhering to the 
following criteria:

- [MAY](#mustmay) contain one or more [op](#op)s at the root of the op collection directory.

**op collection file bundle**  
```TEXT
  |-- .op-collection.yml
  |-- .common
  ... (op bundles)
```

# Default Op Collection
A default op collection may be designated for a directory according to the following rules:

- [MUST](#mustmay) be contained by an `.opspec` directory

In the event a default op collection is not present in a directory, its
nearest ancestor [MUST](#mustmay) be used as the effective default

