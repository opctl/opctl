# Name

[MUST](index.md#mustmay) be named `op.yml`


# Format

[MUST](index.md#mustmay) be
[v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)


# Schema

attributes:

| name        | description           | schema                                                                | required | default            |
|:------------|:----------------------|:----------------------------------------------------------------------|:---------|:-------------------|
| name        | Name of the op        | [FILE_SAFE_NAME](index.md#file_safe_name)                             | true     |                    |
| description | Description of the op | [string](http://yaml.org/type/str.html)                               | true     |                    |
| inputs      | Inputs to the op      | [sequence](http://yaml.org/type/seq.html) of [Parameter](#parameter)s | false    |                    |
| run         | What the op runs      | [Run Statement](#run-statement)                                       | false    | docker-compose.yml |
| version     | Version of the op     | [v2.0.0 SemVer](http://semver.org/spec/v2.0.0.html)                   | false    |                    |


## Run Statement

attributes

| name | description | schema    | required | default |
|:-----|:------------|:----------|:---------|:--------|
| op   | Op to run   | [Op](#op) | true     |         |

or

| name     | description                                 | schema                                                                        | required | default |
|:---------|:--------------------------------------------|:------------------------------------------------------------------------------|:---------|:--------|
| parallel | A sequence of statements to run in parallel | [sequence](http://yaml.org/type/seq.html) of [Run Statement](#run-statement)s | true     |         |

or

| name   | description                              | schema                                                                        | required | default |
|:-------|:-----------------------------------------|:------------------------------------------------------------------------------|:---------|:--------|
| serial | A sequence of statements to run serially | [sequence](http://yaml.org/type/seq.html) of [Run Statement](#run-statement)s | true     |         |


### Op

attributes:

| name | description        | schema                    | required | default |
|:-----|:-------------------|:--------------------------|:---------|:--------|
| ref  | Reference to an op | [OP_REF](index.md#op_ref) | true     |         |


## Parameter

attributes

| name     | description                | schema                                    | required | default |
|:---------|:---------------------------|:------------------------------------------|:---------|:--------|
| type     | Type of parameter          | Name of a defined parameter type          | false    | String  |
| isSecret | If the parameter is secret | [bool](http://yaml.org/type/bool.html)    | false    | false   |
| name     | Name of the parameter      | [FILE_SAFE_NAME](index.md#file_safe_name) | true     |         |
| default  | default value              | [string](http://yaml.org/type/str.html)   | false    |         |
