# Table of Contents

- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Examples](#examples)
- [Bundles](#bundles)
- [Ops](#ops)
    - [Op Definitions](#op-definitions)
- [Collections](#collections)
    - [Collection Definitions](#collection-definitions)
    - [Default Collections](#default-collections)
- [Registry](#registry)
    - [Registry API](#registry-api)
- [Engine](#engine)
    - [Engine API](#engine-api)

# Introduction

## Purpose

Opspec is a specification for defining, distributing, and running ops
(operations).

Primary concerns of opspec are to make operations:

- composable
- portable
- versionable

## MUST/MAY/RECOMMENDED

as defined by [RFC 2119](https://tools.ietf.org/html/rfc2119)

## Examples

Self contained [examples](examples/) are included with this spec.

It is [RECOMMENDED](#mustmayrecommended) integrators use them to
document/demonstrate usage of their integrations.


# Bundles

Bundles are directories containing a manifest and artifacts (dependent
files/folders).

# Ops

An op is a runnable task.

## Op Definitions

Ops are defined via [Bundles](#bundles). Valid op bundles meet the
following criteria:

- [MUST](#mustmayrecommended) contain a
  [collection.yml file](op.yml-file.md) at their root.
- Name [MUST](#mustmayrecommended) match that of the op they contain.


# Collections

One or more [ops](#ops), grouped together physically (via embedding)
and/or logically (via reference).

## Collection Definitions

Collection are defined via [Bundles](#bundles). Valid collection bundles
meet the following criteria:

- [MUST](#mustmayrecommended) contain a
  [collection.yml file](collection.yml-file.md) at their root.
- Name [MUST](#mustmayrecommended) match that of the collection they
  contain unless [designated as default](#default-collections)


# Default Collections

Directories within hierarchical filesystems, [MAY](#mustmayrecommended)
contain [default collections](#default-collections). For a collection to
be designated as default, it's bundle [MUST](#mustmayrecommended) be
named `.opspec`.

[default collections](#default-collections) [MUST](#mustmayrecommended),
by default, be effective within their containing directory.


# Registry

Registries store [bundle](#bundles), enabling centralized publication,
discovery, and consumption.

## Registry API

Registries [MUST](#mustmayrecommended) implement the
[registry api](registry-oai_spec.yaml)


# Engine

Engines run ops.

## Engine API

Engines [MUST](#mustmayrecommended) implement the
[engine api](engine-oai_spec.yaml)
