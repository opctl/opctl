---
sidebar_label: Introduction
title: Introduction
---
Opspec (operational specification) is a declarative language designed for the sole purpose of describing operations (ops).

## Structure
Each op is defined as a directory.

Reference:
- [{OP_DIRECTORY}](reference/structure/op-directory/index)
    - [op.yml](reference/structure/op-directory/index#opyml)
        - [name](reference/structure/op-directory/op/index#name)
        - [description](reference/structure/op-directory/op/index#description)
        - [inputs](reference/structure/op-directory/op/index#inputs)/[outputs](reference/structure/op-directory/op/index#outputs)
            - [{PARAMETER_NAME}](reference/structure/op-directory/op/parameter/index)
                > one of...

                - [array](reference/structure/op-directory/op/parameter/array)
                    - [constraints](reference/structure/op-directory/op/parameter/array#constraints)
                    - [default](reference/structure/op-directory/op/parameter/array#default)
                    - [description](reference/structure/op-directory/op/parameter/array#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/array#issecret)
                - [boolean](reference/structure/op-directory/op/parameter/boolean)
                    - [default](reference/structure/op-directory/op/parameter/boolean#default)
                    - [description](reference/structure/op-directory/op/parameter/boolean#description)
                - [dir](reference/structure/op-directory/op/parameter/dir)
                    - [default](reference/structure/op-directory/op/parameter/dir#default)
                    - [description](reference/structure/op-directory/op/parameter/dir#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/dir#issecret)
                - [file](reference/structure/op-directory/op/parameter/file)
                    - [default](reference/structure/op-directory/op/parameter/file#default)
                    - [description](reference/structure/op-directory/op/parameter/file#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/file#issecret)
                - [number](reference/structure/op-directory/op/parameter/number)
                    - [constraints](reference/structure/op-directory/op/parameter/number#constraints)
                    - [default](reference/structure/op-directory/op/parameter/number#default)
                    - [description](reference/structure/op-directory/op/parameter/number#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/number#issecret)
                - [object](reference/structure/op-directory/op/parameter/object)
                    - [constraints](reference/structure/op-directory/op/parameter/object#constraints)
                    - [default](reference/structure/op-directory/op/parameter/object#default)
                    - [description](reference/structure/op-directory/op/parameter/object#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/object#issecret)
                - [socket](reference/structure/op-directory/op/parameter/socket)
                    - [description](reference/structure/op-directory/op/parameter/socket#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/socket#issecret)
                - [string](reference/structure/op-directory/op/parameter/string)
                    - [constraints](reference/structure/op-directory/op/parameter/string#constraints)
                    - [default](reference/structure/op-directory/op/parameter/string#default)
                    - [description](reference/structure/op-directory/op/parameter/string#description)
                    - [isSecret](reference/structure/op-directory/op/parameter/string#issecret)
        - [opspec](reference/structure/op-directory/op/index#opspec)
        - [run](reference/structure/op-directory/op/index#run)
            > one of...

            - [container](reference/structure/op-directory/op/call/container/index)
                - [cmd](reference/structure/op-directory/op/call/container/index#cmd)
                - [dirs](reference/structure/op-directory/op/call/container/index#dirs)
                - [envVars](reference/structure/op-directory/op/call/container/index#envvars)
                - [files](reference/structure/op-directory/op/call/container/index#files)
                - [image](reference/structure/op-directory/op/call/container/image/index)
                    - [ref](reference/structure/op-directory/op/call/container/image#ref)
                    - [pullCreds](reference/structure/op-directory/op/call/container/image#pullcreds)
                - [name](reference/structure/op-directory/op/call/container/index#name)
                - [ports](reference/structure/op-directory/op/call/container/index#ports)
                - [sockets](reference/structure/op-directory/op/call/container/index#sockets)
                - [workDir](reference/structure/op-directory/op/call/container/index#workdir)
            - [op](reference/structure/op-directory/op/call/op)
                - [inputs](reference/structure/op-directory/op/call/op/index#inputs)
                - [outputs](reference/structure/op-directory/op/call/op/index#outputs)
                - [pullCreds](reference/structure/op-directory/op/call/op/index#pullcreds)
                - [ref](reference/structure/op-directory/op/call/op/index#ref)
            - [parallel](reference/structure/op-directory/op/call/index#parallel)
            - [parallelLoop](reference/structure/op-directory/op/call/parallelloop)
                - [range](reference/structure/op-directory/op/call/parallelloop#range)
                - [run](reference/structure/op-directory/op/call/parallelloop#run)
                - [vars](reference/structure/op-directory/op/call/parallelloop#vars)
            - [serial](reference/structure/op-directory/op/call/index#serial)
            - [serialLoop](reference/structure/op-directory/op/call/index#serialloop)
                - [range](reference/structure/op-directory/op/call/serialloop#range)
                - [run](reference/structure/op-directory/op/call/serialloop#run)
                - [until](reference/structure/op-directory/op/call/serialloop#until)
                - [vars](reference/structure/op-directory/op/call/serialloop#vars)
        - [version](reference/structure/op-directory/op/index#version)
    - [icon.svg](reference/structure/op-directory/index#iconsvg)

## Types
Values in opspec are typed.

Reference:
- [array](reference/types/array)
- [boolean](reference/types/boolean)
- [dir](reference/types/dir)
- [file](reference/types/file)
- [number](reference/types/number)
- [object](reference/types/object)
- [socket](reference/types/socket)
- [string](reference/types/string)

## Scoping
Variables in opspec are scoped to each operation. [Parameters](../op/parameter/index) allow passing values into or out of this scope.