the dir contains opspec packages used to
[fuzz test](https://en.wikipedia.org/wiki/Fuzzing) package validation.

The package naming convention used is as follows:

(valid|invalid)_(objectPath)_(scenario)

examples:

| name                 | objectPath | scenario |
|:---------------------|:-----------|:---------|
| invalid__yml         | -          | yml      |
| invalid_inputs_type  | inputs     | type     |
| invalid_outputs_type | outputs    | type     |
| valid_all            | -          | all      |
