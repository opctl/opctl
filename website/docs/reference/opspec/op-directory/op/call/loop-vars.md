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
A [variable reference [string]](../variable-reference.md) each iterations index will be bound to.

### key
A [variable reference [string]](../variable-reference.md) each iterations key will be bound to.

Behavior varies based on the range value:  

|range|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item index|
|object|Variable set to current property name|

### value
A [variable reference [string]](../variable-reference.md) each iterations value will be bound to.

Behavior varies based on the range value:  

|range value|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item|
|object|Variable set to current property value|