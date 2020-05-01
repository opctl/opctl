---
title: Visualization
---

## Local Op
1. Open a terminal and `cd` to any directory that contains ops somewhere in it's subtree.
1. Use the [opctl ui subcommand](../reference/ui.md#ui) to open the web UI to the current directory.
   ```sh
   opctl ui
   ```
1. Use the explorer panel to open and visualize an `op.yml`.
1. Try panning and zooming.

## Remote Op
1. Open a terminal and `cd` to a location that has some ops.
1. Use the [opctl ui subcommand](../reference/ui.md#ui) to open the web UI to the remote git tag [github.com/opspec-pkgs/_.op.create#3.3.1](https://github.com/opspec-pkgs/_.op.create/tree/3.3.1).
   ```sh
   opctl ui github.com/opspec-pkgs/_.op.create#3.3.1
   ```
1. Use the explorer panel to open and visualize the `op.yml`.
   > If any private remote refs are encountered and aren't cached, you'll be prompted for pullCreds.
1. Try panning and zooming.