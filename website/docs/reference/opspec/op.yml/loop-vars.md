---
title: Loop variables
---

Loop variables is an object declaring variables that hold values from a [rangable value](rangeable-value).

## Properties

Each property maps a component of the rangeable value to a [variable reference](variable-reference).

### `index`

A [variable reference](variable-reference.md) to store the integer index of the iteration in. Indexes are available for both array and object ranges.

### `key`

A [variable reference](variable-reference.md) to store the property name of each item in the object being iterated over. This is not supported for arrays.

### `value`

A [variable reference](variable-reference.md) to store the value of each item in the array or object being iterated over. Values are supported for both array and object ranges.
