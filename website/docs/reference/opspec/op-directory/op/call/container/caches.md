---
title: Caches [object]
---

An object mapping absolute container paths to cache names.

```yaml
run:
  container:
    image: { ref: golang:1.23 }
    caches:
      /root/.cache/go-build: go-build
      /go/pkg/mod: gomod
```

Each value must be a string literal or [string expression](../../../../types/string.md#initialization), such as `npm-$(lockHash)`, evaluated in the call's scope. The resulting cache name must be non-empty and may contain only ASCII letters, digits, `.`, `_`, and `-`. The names `.` and `..` are invalid. A container path cannot appear in both `dirs` and `caches`.

Opspec has no hash function. Compute a lockfile hash in a prior call, then use it in the cache name. This example writes the hash to a file bound to `lockHash`, which the next call reads through string interpolation. The example assumes `op.yml` lives at `.opspec/npm-install/` in the project. Remove the trailing newline so it does not become part of the name.

```yaml
name: npm-install
run:
  serial:
    - container:
        image: { ref: alpine:3.20 }
        files:
          /package-lock.json: $(../../package-lock.json)
          /lockHash: $(lockHash)
        cmd:
          - sh
          - -ec
          - sha256sum /package-lock.json | cut -d ' ' -f 1 | tr -d '\n' > /lockHash
    - container:
        image: { ref: node:22 }
        dirs: { /app: $(../..) }
        workDir: /app
        caches:
          /root/.npm: npm-$(lockHash)
        cmd: [npm, ci]
```

The node stores each cache dir at `<data-dir>/caches/<name>`. The cache root is `<data-dir>/caches`. Contents persist between runs on the node, but are never guaranteed. An empty or partial cache dir must always be a valid starting point. Failed or killed calls keep whatever they wrote.

Cache names are global per node. Calls using the same name share contents across ops, even when they mount the cache dir at different container paths.

Concurrent calls share the directory without locking. Tools that are not safe to run concurrently, such as apt, may corrupt their cache contents.

Deleting a cache dir or the cache root is always safe when no call is using it. Later calls recreate the directories as needed.
