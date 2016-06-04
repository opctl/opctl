# Table of Contents
- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Requirements](#requirements)
    - [Terminology](#terminology)
- [Ops](#ops)
- [Op Collections](#op-collections)

# Introduction

## Purpose
Op Spec is an op specification format.

Primary concerns of op spec are to make operations:
- portable across runtimes
- distributable
- discoverable
- versionable

## Requirements
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", 
"RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in 
[RFC 2119](https://tools.ietf.org/html/rfc2119).

An implementation is not compliant if it fails to satisfy one or more of the MUST or REQUIRED 
level requirements for the protocols it implements. An implementation that satisfies all the MUST or 
REQUIRED level and all the SHOULD level requirements for its protocols is said to be "unconditionally 
compliant"; one that satisfies all the MUST level requirements but not all the SHOULD level 
requirements for its protocols is said to be "conditionally compliant."

## Terminology

### FILE_SAFE_NAME
a string matching the regex `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`

### FILE_BUNDLE
a file directory with a defined structure (see <a href="https://en.wikipedia.org/wiki/Bundle_(OS_X)">Bundle_(OS_X)</a> for similar usage)


# Ops
Ops are processes or tasks. Ops can be composed, i.e. ops can consist of other constituent ops. 

## Op Bundle
Ops [MUST](./index.md#requirements) be stored as op bundles (see [FILE_BUNDLE](#file_bundle)).

**Manifest**
An [op manifest file](op-manifest-file.md) [MUST](./index.md#requirements) 
exist at the root of an op bundle.

**Docker Composition**
An [op docker compose file](./op-docker-compose-file.md) [MUST](./index.md#requirements) 
exist at the root of an op bundle.

**Tree**
```TEXT
  |-- op.yml
  |-- docker-compose.yml
  ... (op specific files/dirs)
```

# Op Collections
Op collections consist of one or more [ops](#ops)

## Op Collection Bundle
Op collections [MUST](./index.md#requirements) be stored as op collection bundles (see [FILE_BUNDLE](#file_bundle)).

**Manifest**
An [op collection manifest file](op-collection-manifest-file.md)
[MUST](./index.md#requirements) exist at the root of an op collection bundle.

**Embedded Op Bundles**
One or more [op bundle](#op-bundle)s [MAY](./index.md#requirements) be embedded
by including them as child directories.

**Default Designation** 
An [op collection bundle](#op-collection-bundle) who's containing folder is named `.opspec` 
[MUST](./index.md#requirements) be considered the default op collection at that path.

**Tree**
```TEXT
  |-- .op-collection.yml
  |-- .common
  ... (op bundles)
```
