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
A [call](../call/index.md) object defining a call run on each iteration of the loop

### until
An array of [predicates](predicate.md) which must be true for the loop to exit.

### vars
A [loop-vars](loop-vars.md) object binding iteration info to variables.