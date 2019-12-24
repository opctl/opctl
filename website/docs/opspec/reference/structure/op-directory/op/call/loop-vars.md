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
A string defining the name of a variable to set equal to the current loop index.

### key
A string defining the name of a variable to set equal to the current loop key.

Behavior varies based on the range value:  

|range|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item index|
|object|Variable set to current property name|

### value
A string defining the name of a variable to set equal to the current loop value.

Behavior varies based on the range value:  

|range value|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item|
|object|Variable set to current property value|