---
sidebar_label: Introduction
title: Introduction
---
Opspec (operational specification) is a declarative language designed for the sole purpose of describing operations (ops).

## Structure
Each op is defined as a directory.

Reference:
- [{OP_DIRECTORY}](reference/structure/op-directory/index.md)
    - [op.yml](reference/structure/op-directory/index.md#opyml)
        - [name](reference/structure/op-directory/op/index.md#name)
        - [description](reference/structure/op-directory/op/index.md#description)
        - [inputs](reference/structure/op-directory/op/index.md#inputs)/[outputs](reference/structure/op-directory/op/index.md#outputs)
            - [{PARAMETER_NAME}](reference/structure/op-directory/op/parameter/index.md)
                > one of...

                - [array](reference/structure/op-directory/op/parameter/array.md)
                    - [constraints](reference/structure/op-directory/op/parameter/array.md#constraints)
                    - [default](reference/structure/op-directory/op/parameter/array.md#default)
                    - [description](reference/structure/op-directory/op/parameter/array.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/array.md#issecret)
                - [boolean](reference/structure/op-directory/op/parameter/boolean.md)
                    - [default](reference/structure/op-directory/op/parameter/boolean.md#default)
                    - [description](reference/structure/op-directory/op/parameter/boolean.md#description)
                - [dir](reference/structure/op-directory/op/parameter/dir.md)
                    - [default](reference/structure/op-directory/op/parameter/dir.md#default)
                    - [description](reference/structure/op-directory/op/parameter/dir.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/dir.md#issecret)
                - [file](reference/structure/op-directory/op/parameter/file.md)
                    - [default](reference/structure/op-directory/op/parameter/file.md#default)
                    - [description](reference/structure/op-directory/op/parameter/file.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/file.md#issecret)
                - [number](reference/structure/op-directory/op/parameter/number.md)
                    - [constraints](reference/structure/op-directory/op/parameter/number.md#constraints)
                    - [default](reference/structure/op-directory/op/parameter/number.md#default)
                    - [description](reference/structure/op-directory/op/parameter/number.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/number.md#issecret)
                - [object](reference/structure/op-directory/op/parameter/object.md)
                    - [constraints](reference/structure/op-directory/op/parameter/object.md#constraints)
                    - [default](reference/structure/op-directory/op/parameter/object.md#default)
                    - [description](reference/structure/op-directory/op/parameter/object.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/object.md#issecret)
                - [socket](reference/structure/op-directory/op/parameter/socket.md)
                    - [description](reference/structure/op-directory/op/parameter/socket.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/socket.md#issecret)
                - [string](reference/structure/op-directory/op/parameter/string.md)
                    - [constraints](reference/structure/op-directory/op/parameter/string.md#constraints)
                    - [default](reference/structure/op-directory/op/parameter/string.md#default)
                    - [description](reference/structure/op-directory/op/parameter/string.md#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/string.md#issecret)
        - [opspec](reference/structure/op-directory/op/index.md#opspec)
        - [run](reference/structure/op-directory/op/index.md#run)
            > one of...

            - [container](reference/structure/op-directory/op/call/container/index.md)
                - [cmd](reference/structure/op-directory/op/call/container/index.md#cmd)
                - [dirs](reference/structure/op-directory/op/call/container/index.md#dirs)
                - [envVars](reference/structure/op-directory/op/call/container/index.md#envvars)
                - [files](reference/structure/op-directory/op/call/container/index.md#files)
                - [image](reference/structure/op-directory/op/call/container/image/index.md)
                    - [ref](reference/structure/op-directory/op/call/container/image.md#ref)
                    - [pullCreds](reference/structure/op-directory/op/call/container/image.md#pullcreds)
                - [name](reference/structure/op-directory/op/call/container/index.md#name)
                - [ports](reference/structure/op-directory/op/call/container/index.md#ports)
                - [sockets](reference/structure/op-directory/op/call/container/index.md#sockets)
                - [workDir](reference/structure/op-directory/op/call/container/index.md#workdir)
            - [op](reference/structure/op-directory/op/call/op.md)
                - [inputs](reference/structure/op-directory/op/call/op.md#inputs)
                - [outputs](reference/structure/op-directory/op/call/op.md#outputs)
                - [pullCreds](reference/structure/op-directory/op/call/op.md#pullcreds)
                - [ref](reference/structure/op-directory/op/call/op.md#ref)
            - [parallel](reference/structure/op-directory/op/call/index.md#parallel)
            - [parallelLoop](reference/structure/op-directory/op/call/parallel-loop.md)
                - [range](reference/structure/op-directory/op/call/parallel-loop.md#range)
                - [run](reference/structure/op-directory/op/call/parallel-loop.md#run)
                - [vars](reference/structure/op-directory/op/call/parallel-loop.md#vars)
            - [serial](reference/structure/op-directory/op/call/index.md#serial)
            - [serialLoop](reference/structure/op-directory/op/call/serial-loop.md)
                - [range](reference/structure/op-directory/op/call/serial-loop.md#range)
                - [run](reference/structure/op-directory/op/call/serial-loop.md#run)
                - [until](reference/structure/op-directory/op/call/serial-loop.md#until)
                - [vars](reference/structure/op-directory/op/call/serial-loop.md#vars)
        - [version](reference/structure/op-directory/op/index.md#version)
    - [icon.svg](reference/structure/op-directory/index.md#iconsvg)

## Types
Values in opspec are typed.

Reference:
- [array](reference/types/array.md)
- [boolean](reference/types/boolean.md)
- [dir](reference/types/dir.md)
- [file](reference/types/file.md)
- [number](reference/types/number.md)
- [object](reference/types/object.md)
- [socket](reference/types/socket.md)
- [string](reference/types/string.md)

## Scoping
Variables in opspec are scoped to each operation. [Parameters](../op/parameter/index.md) allow passing values into or out of this scope.