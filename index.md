# Table of Contents
- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Terminology](#terminology)
- [Op](#op)
    - [Op Bundle](#op-bundle)
- [Op Collection](#op-collection)
    - [Op Collection Bundle](#op-collection-bundle)

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
a file directory with a defined structure allowing related files to be grouped together as a conceptually single item (see <a href="https://en.wikipedia.org/wiki/Bundle_(OS_X)">Bundle_(OS_X)</a> for similar usage)

### MUST/MAY
as defined by [RFC 2119](https://tools.ietf.org/html/rfc2119)

# Op
An op is a process or task. Ops can be composed, i.e. ops can consist of other constituent ops. 

## Op Bundle
Ops [MUST](#mustmay) be stored as op bundles (see [FILE_BUNDLE](#file_bundle)).

**Manifest**  
An [op manifest file](op-manifest-file.md) [MUST](#mustmay) 
exist at the root of an op bundle.

**Docker Composition**  
An [op docker compose file](./op-docker-compose-file.md) [MUST](#mustmay) 
exist at the root of an op bundle.

**Tree**  
```TEXT
  |-- op.yml
  |-- docker-compose.yml
  ... (op specific files/dirs)
```

# Op Collection
An op collection consists of one or more [op](#op)s

## Op Collection Bundle
Op collections [MUST](#mustmay) be stored as op collection bundles (see [FILE_BUNDLE](#file_bundle)).

**Embedded Op Bundles**  
One or more [op bundle](#op-bundle)s [MAY](#mustmay) be embedded
by including them as child directories.

**Default Designation**  
An [op collection bundle](#op-collection-bundle) who's containing folder is named `.opspec` 
[MUST](#mustmay) be considered the default op collection at that path.

**Tree**  
```TEXT
  |-- .op-collection.yml
  |-- .common
  ... (op bundles)
```
