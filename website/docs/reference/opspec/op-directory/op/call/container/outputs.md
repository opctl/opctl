---
sidebar_label: Outputs file
title: Container outputs file
---

A nonempty `container.outputs` map requests an outputs file. Each key is a JSON member name and each value is a whole scope write reference. Keys and target names contain letters, digits, underscores, or hyphens.

```yaml
container:
  image: {ref: alpine:3}
  outputs:
    version: $(buildVersion)
    count: $(buildCount)
  cmd:
    - sh
    - -ce
    - |
      printf '%s' '{"version":"v1.2.3","count":3}' > "$OPCTL_OUTPUTS"
```

opctl mounts a fresh writable directory at `/opctl/outputs` and sets `OPCTL_OUTPUTS=/opctl/outputs/values.json`. The file starts absent. Applications may write it directly or write a temporary sibling and rename it. Authored environment variables named `OPCTL_OUTPUTS` and mounts that overlap `/opctl/outputs` are rejected only when outputs are requested. Calls without bindings get no new environment entry or mount.

Write one UTF-8 JSON object, at most 1 MiB including whitespace and unrequested members. Strings remain exact, including newlines and literal `$()` text. Numbers use finite float64 values; use strings for large exact identifiers. Booleans, arrays, and objects retain native JSON data. Paths in JSON are strings, not file handles.

Every requested key must exist. Extra keys are validated and ignored. Empty files, null anywhere, duplicate keys at any depth, nesting beyond 100 containers, non-object roots, trailing documents or garbage, symlinks, FIFOs, and non-regular files fail the call. Diagnostics omit file contents. Collection and binding happen after successful execution, and any failure publishes no new results. Runtime failures take precedence over file errors.

Secret output declarations and sensitive targets keep values out of completion events while later calls receive them privately. Applications must avoid printing secrets themselves. The file lives under the private call directory in the node data dir. opctl removes only its outputs directory on completion or failure; a node crash can leave scratch data. The 1 MiB read limit is not a disk write quota.

Docker supports outputs files, including non-root writers. Embedded requires a same-path Lima share and a macOS mount/UID smoke before claiming support. Outputs-file-only YAML on Kubernetes uses the existing shared PVC under `/root/opctl`; paths outside it are rejected. A live shared-PVC smoke remains required. Kubernetes stream bindings and Go bootstrap remain unsupported.

## Typed members

Add `outputParams` beside `outputs` to validate produced JSON before binding it into scope:

```yaml
container:
  image: {ref: alpine:3}
  outputs: {count: $(buildCount), metadata: $(metadata)}
  outputParams:
    count: {number: {constraints: {minimum: 1, format: integer}}}
    metadata:
      object:
        isSecret: true
        constraints:
          properties: {revision: {type: string}}
          required: [revision]
  cmd: [sh, -ce, 'printf ''%s'' ''{"count":3,"metadata":{"revision":"abc123"}}'' > "$OPCTL_OUTPUTS"']
```

Declaration keys name literal JSON members, so `count` validates the member before `outputs` renames it to `buildCount`. Every declared member must exist, including members without a binding. A missing member, wrong JSON kind, or failed constraint fails the container and publishes no new outputs. A JSON string `"3"` fails a Number declaration even though ordinary op parameters can coerce it.

Declarations use the existing parameter grammar for String, Number, Boolean, Array, and Object, with descriptions, constraints, and supported `isSecret` flags. Defaults and Dir/File/Socket declarations are rejected. Use existing mount and socket bindings for those kinds. A declaration map must be nonempty and requires nonempty `outputs`. Bindings without declarations keep their existing behavior.

A declared secret protects its bound destination before execution. Public declarations cannot remove inherited protection. Secret constraint errors omit produced values. Captured stdout and stderr remain String values, independent of `outputParams`; a declaration named `stdout` describes a JSON member, not the process stream. Applications still control what they print to public logs.

Go graphs using typed container declarations require opspec schema 4. Older nodes fail with an upgrade diagnostic.
