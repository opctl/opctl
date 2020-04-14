---
title: Loop Vars [object]
---

An object naming variables to assign each iterations info to.

## Properties
- may have
  - [index](#index)
  - [key](#key)
  - [value](#value)

### index
An [identifier [string]](../identifier.md) naming a variable to bind each iterations index to.

### key
An [identifier [string]](../identifier.md) naming a variable to bind each iterations key to.

Behavior varies based on the range value:  

|range|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item index|
|object|Variable set to current property name|

### value
An [identifier [string]](../identifier.md) naming a variable to bind each iterations value to.

Behavior varies based on the range value:  

|range value|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item|
|object|Variable set to current property value|