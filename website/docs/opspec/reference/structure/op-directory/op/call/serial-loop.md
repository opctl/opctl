---
title: Serial Loop Call [object]
---

An object defining a call loop in which each iteration happens in serial (one after another in order)

## Properties:
- may have
  - [run](#run)
  - [vars](#vars)
- must have at least one of
  - [range](#range)
  - [until](#until)

### range
A [rangeable value](rangeable-value.md) to loop over.

### run
A [call [object]](../call/index.md) defining a call run on each iteration of the loop

### until
An array of [predicate [object]](predicate.md)s which must be true for the loop to exit.

### vars
A [loop-vars [object]](loop-vars.md) binding iteration info to variables.