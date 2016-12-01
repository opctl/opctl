# Table of Contents

- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Examples](#examples)
- [Bundles](#bundles)
- [Ops](#ops)
    - [Op Definitions](#op-definitions)
- [Collections](#collections)
    - [Collection Definitions](#collection-definitions)
    - [Designated Collections](#designated-collections)
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

## MUST/MAY

as defined by [RFC 2119](https://tools.ietf.org/html/rfc2119)

## Examples

Self contained [examples](examples/) are included with this spec.

It is RECOMMENDED integrators use them to document/demonstrate usage of
their integrations.


# Bundles

Bundles are directories containing a manifest and artifacts (dependent
files/folders).

# Ops

An op is a runnable task.

## Op Definitions

Ops are defined via a [Bundle](#bundles). Op bundles [MUST](#mustmay)
contain an [op.yml file](op.yml-file.md) at their root.


# Collections

One or more [op](#ops), grouped together physically (via embedding)
and/or logically (via reference).

## Collection Definitions

Collections are defined via a [Bundle](#bundles). Collection bundles
[MUST](#mustmay) contain a [collection.yml file](collection.yml-file.md)
at their root.


# Designated Collections

Directories within hierarchical filesystems, [MAY](#mustmay) contain a
collection bundle named `.opspec` which if present, [MUST](#mustmay) be
treated as their designated collection.


# Registry

Registries store [bundle](#bundles), enabling centralized publication,
discovery, and consumption.

## Registry API

Registries [MUST](#mustmay) implement the
[registry api](registry-oai_spec.yaml)


# Engine

Engines run ops.

## Engine API

Engines [MUST](#mustmay) implement the
[engine api](engine-oai_spec.yaml)
