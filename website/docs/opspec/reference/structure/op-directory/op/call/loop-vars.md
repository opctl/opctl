---
title: Loop Vars [object]
---

An object binding variables to iteration info (index, key, value) assigned to each iteration.

## Properties
- may have
  - [index](#index)
  - [key](#key)
  - [value](#value)

### index
A [variable-reference [string]](../variable-reference.md) designating a variable to bind each iterations index to.

### key
A [variable-reference [string]](../variable-reference.md) designating a variable to bind each iterations key to.

Behavior varies based on the range value:  

|range|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item index|
|object|Variable set to current property name|

### value
A [variable-reference [string]](../variable-reference.md) designating a variable to bind each iterations value to.

Behavior varies based on the range value:  

|range value|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item|
|object|Variable set to current property value|