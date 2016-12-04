# Table of Contents

- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Examples](#examples)
- [Bundles](#bundles)
- [Ops](#ops)
    - [Op Bundles](#op-bundles)
        - [op.yml file](#opyml-file)
- [Collections](#collections)
    - [Collection Bundles](#collection-bundles)
        - [collection.yml file](#collectionyml-file)
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

Bundles are directories containing a manifest (at their root) and
optional artifacts (dependent files/folders).

# Ops

A runnable task.

## Op Bundles

Ops are defined via [Bundles](#bundles). Valid op bundles meet the
following criteria:

- [MUST](#mustmayrecommended) contain an [op.yml file](#opyml-file) at
  their root.
- Name [MUST](#mustmayrecommended) match that of the op they contain.

## op.yml file

`op.yml` files are the manifest for op bundles. Valid `op.yml`
files meet the following criteria:

- [MUST](index.md#mustmayrecommended) be named `op.yml`
- [MUST](index.md#mustmayrecommended) be
  [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- [MUST](index.md#mustmayrecommended) validate against
  [schema/manifest.json#definitions/opManifest](schema/manifest.json#definitions/opManifest)


# Collections

One or more [ops](#ops), grouped together physically (via embedding)
and/or logically (via reference).

## Collection Bundles

Collection are defined via [Bundles](#bundles). Valid collection bundles
meet the following criteria:

- [MUST](#mustmayrecommended) contain a
  [collection.yml file](#collectionyml-file) at their root.
- Name [MUST](#mustmayrecommended) match that of the collection they
  contain unless [designated as default](#default-collections)

## collection.yml file

`collection.yml` files are the manifest for collection bundles. Valid
`collection.yml` files meet the following criteria:

- [MUST](index.md#mustmayrecommended) be named `collection.yml`
- [MUST](index.md#mustmayrecommended) be
  [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- [MUST](index.md#mustmayrecommended) validate against
  [schema/manifest.json#definitions/collectionManifest](schema/manifest.json#definitions/collectionManifest)


## Default Collections

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
[registry api](registry-oai_spec.yml)


# Engine

Engines run ops.

## Engine API

Engines [MUST](#mustmayrecommended) implement the
[engine api](engine-oai_spec.yml)

