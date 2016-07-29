# Table of Contents

- [Introduction](#introduction)
    - [Purpose](#purpose)
    - [Terminology](#terminology)
- [Bundle](#bundle)
- [Op](#op)
- [Collection](#collection)
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
- portable
- versionable

## Terminology

### FILE_SAFE_NAME

a string matching the regex `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`

### MUST/MAY

as defined by [RFC 2119](https://tools.ietf.org/html/rfc2119)

# Bundle

Bundles are directories containing a manifest and artifacts (dependent
files/folders).

Once defined, **bundles**, can be distributed via any traditional means
of file transfer or published to a [Registry](#registry)

# Op

A task; work.

Ops are defined via a [Bundle](#bundle) containing an
[op.yml file](op.yml-file.md) and optionally a
[docker-compose.yml file](docker-compose.yml.md).


# Collection

One or more [op](#op)s, grouped together physically (via embedding)
and/or logically (via reference).

Collections are defined via a [Bundle](#bundle) containing a
[collection.yml file](collection.yml-file.md).


## Default Collection

A default collection may be designated for a directory according to the
following criteria:

- [MUST](#mustmay) be contained by an `.opspec` directory

In the event a default collection is not present in a directory, its
nearest ancestor [MUST](#mustmay) be used as the effective default


# Registry

Registries store [bundle](#bundle)s, enabling centralized publication,
discovery, and consumption.

## Registry API

Registries [MUST](#mustmay) implement the
[registry api](registry-oai_spec.yaml)


# Engine

Engines run ops.

## Engine API

Engines [MUST](#mustmay) implement the
[engine api](engine-oai_spec.yaml)
