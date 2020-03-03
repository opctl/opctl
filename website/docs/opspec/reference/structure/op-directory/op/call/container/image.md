---
title: Image [object]
---

An object which defines the image of a container call.

## Properties
- must have
  - [ref](#ref)
- may have
  - [pullCreds](#pullcreds)

### ref
A [string initializer](../../../../../types/string.md#initialization) referencing a container image.

### Example ref ([official ubuntu 19.10](https://hub.docker.com/_/ubuntu))
`ref: 'ubuntu:19.10'`

### pullCreds
A [pull-creds](../pull-creds.md) object defining creds used to pull the image from a private source.