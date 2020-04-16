---
sidebar_label: Index
title: Opspec
---
Opspec (portmanteau of operation specification) is a language designed to portably and fully define operations (ops).

## Structure
Each op is defined as a directory.

Reference:
- [{OP_DIRECTORY}](op-directory/index.md)
    - [op.yml](op-directory/index.md#opyml)
        - [name](op-directory/op/index.md#name)
        - [description](op-directory/op/index.md#description)
        - [inputs](op-directory/op/index.md#inputs)/[outputs](op-directory/op/index.md#outputs)
            - [{PARAMETER_NAME}](op-directory/op/parameter/index.md)
                > one of...

                - [array](op-directory/op/parameter/array.md)
                    - [constraints](op-directory/op/parameter/array.md#constraints)
                    - [default](op-directory/op/parameter/array.md#default)
                    - [description](op-directory/op/parameter/array.md#description)
                    - [isSecret](op-directory/op/parameter/array.md#issecret)
                - [boolean](op-directory/op/parameter/boolean.md)
                    - [default](op-directory/op/parameter/boolean.md#default)
                    - [description](op-directory/op/parameter/boolean.md#description)
                - [dir](op-directory/op/parameter/dir.md)
                    - [default](op-directory/op/parameter/dir.md#default)
                    - [description](op-directory/op/parameter/dir.md#description)
                    - [isSecret](op-directory/op/parameter/dir.md#issecret)
                - [file](op-directory/op/parameter/file.md)
                    - [default](op-directory/op/parameter/file.md#default)
                    - [description](op-directory/op/parameter/file.md#description)
                    - [isSecret](op-directory/op/parameter/file.md#issecret)
                - [number](op-directory/op/parameter/number.md)
                    - [constraints](op-directory/op/parameter/number.md#constraints)
                    - [default](op-directory/op/parameter/number.md#default)
                    - [description](op-directory/op/parameter/number.md#description)
                    - [isSecret](op-directory/op/parameter/number.md#issecret)
                - [object](op-directory/op/parameter/object.md)
                    - [constraints](op-directory/op/parameter/object.md#constraints)
                    - [default](op-directory/op/parameter/object.md#default)
                    - [description](op-directory/op/parameter/object.md#description)
                    - [isSecret](op-directory/op/parameter/object.md#issecret)
                - [socket](op-directory/op/parameter/socket.md)
                    - [description](op-directory/op/parameter/socket.md#description)
                    - [isSecret](op-directory/op/parameter/socket.md#issecret)
                - [string](op-directory/op/parameter/string.md)
                    - [constraints](op-directory/op/parameter/string.md#constraints)
                    - [default](op-directory/op/parameter/string.md#default)
                    - [description](op-directory/op/parameter/string.md#description)
                    - [isSecret](op-directory/op/parameter/string.md#issecret)
        - [opspec](op-directory/op/index.md#opspec)
        - [run](op-directory/op/index.md#run)
            - [if](op-directory/op/call/index.md#if)
            - [name](op-directory/op/call/index.md#name)
            - [needs](op-directory/op/call/index.md#needs)
            > one of...

            - [container](op-directory/op/call/container/index.md)
                - [cmd](op-directory/op/call/container/index.md#cmd)
                - [dirs](op-directory/op/call/container/index.md#dirs)
                - [envVars](op-directory/op/call/container/index.md#envvars)
                - [files](op-directory/op/call/container/index.md#files)
                - [image](op-directory/op/call/container/image/index.md)
                    - [ref](op-directory/op/call/container/image.md#ref)
                    - [pullCreds](op-directory/op/call/container/image.md#pullcreds)
                - [name](op-directory/op/call/container/index.md#name)
                - [ports](op-directory/op/call/container/index.md#ports)
                - [sockets](op-directory/op/call/container/index.md#sockets)
                - [workDir](op-directory/op/call/container/index.md#workdir)
            - [op](op-directory/op/call/op.md)
                - [inputs](op-directory/op/call/op.md#inputs)
                - [outputs](op-directory/op/call/op.md#outputs)
                - [pullCreds](op-directory/op/call/op.md#pullcreds)
                - [ref](op-directory/op/call/op.md#ref)
            - [parallel](op-directory/op/call/index.md#parallel)
            - [parallelLoop](op-directory/op/call/parallel-loop.md)
                - [range](op-directory/op/call/parallel-loop.md#range)
                - [run](op-directory/op/call/parallel-loop.md#run)
                - [vars](op-directory/op/call/parallel-loop.md#vars)
            - [serial](op-directory/op/call/index.md#serial)
            - [serialLoop](op-directory/op/call/serial-loop.md)
                - [range](op-directory/op/call/serial-loop.md#range)
                - [run](op-directory/op/call/serial-loop.md#run)
                - [until](op-directory/op/call/serial-loop.md#until)
                - [vars](op-directory/op/call/serial-loop.md#vars)
        - [version](op-directory/op/index.md#version)
    - [icon.svg](op-directory/index.md#iconsvg)

## Types
Values in opspec are typed.

Reference:
- [array](types/array.md)
- [boolean](types/boolean.md)
- [dir](types/dir.md)
- [file](types/file.md)
- [number](types/number.md)
- [object](types/object.md)
- [socket](types/socket.md)
- [string](types/string.md)

## Scoping
Variables in opspec are scoped to each operation. [Parameters](../op/parameter/index.md) allow passing values into or out of this scope.