---
title: How do I visualize an op?
---

## TLDR;
The opctl UI allows visualizing ops whether defined locally on the opctl host or remotely in a git repo. 

If the op is defined remotely in a git repo and the git repo requires authentication, opctl will prompt for a username & password. 

## Local Op Example
1. Open a terminal and `cd` to any directory that contains ops somewhere in it's subtree.
1. Use the [opctl ui subcommand](../../reference/ui.md#mount) to open the web UI to the current directory.
   ```sh
   opctl ui
   ```
1. Use the explorer panel to open and visualize an `op.yml`.
1. Try panning and zooming.

## Remote Op Example
1. Open a terminal.
1. Use the [opctl ui subcommand](../../reference/ui.md#mount) to open the web UI to the remote git tag [github.com/opspec-pkgs/_.op.create#3.3.1](https://github.com/opspec-pkgs/_.op.create/tree/3.3.1).
   ```sh
   opctl ui github.com/opspec-pkgs/_.op.create#3.3.1
   ```
1. Use the explorer panel to open and visualize the `op.yml`.
   > If any private remote refs are encountered and aren't cached, you'll be prompted for pullCreds.
1. Try panning and zooming.
