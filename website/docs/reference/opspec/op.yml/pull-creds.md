---
title: Pull credentials
---

Pull credentials is an object that defines credentials for a remote data source. This can be used, for example, to authenticate with a private github repository.

## Properties

### `username`

_required_

A [string](../types/string.md#initialization) used as a username.

### `password`

_required_

A [string](../types/string.md#initialization) used as the password.

Certian cli operations will prompt the user for credentials if not provided. Passwords are treated as a secret input when collected like this.
