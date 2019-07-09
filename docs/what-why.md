---
title: What & Why
sidebar_label: What & Why
---

## What are ops
* we think of ops as operations as code
* an op accepts inputs, produces outputs, and may have side effects
* ops are designed to be:
    1. Composable: ops can be composed of smaller ops that are defined to run in serial or parallel
    2. Portable: an op's definition contains everything it needs to run and what inputs it expects, and ops leverage docker containers to run anywhere
    3. Distributable: ops can be referenced remotely, and can be remotely invoked
    4. Versionable: an op is defined in a simple `yaml` file which makes versioning easy using standard source control

## What problems do ops solve?
* automating manual technical operations
* reliable and easy local development for software services
* portable pipelines that live and change with the code they build and deploy
* turning tacit operational knowledge into executable documentation
* providing microservice development teams with an easy to understand, standard interface for operations

## Implementation Goals

- decentralized
- vendor & platform agnostic
- single executable