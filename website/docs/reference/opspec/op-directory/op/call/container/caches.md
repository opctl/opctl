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

Each value must be a non-empty string literal. Cache names may contain only ASCII letters, digits, `.`, `_`, and `-`. The names `.` and `..` are invalid. A container path cannot appear in both `dirs` and `caches`.

The node stores each cache dir at `<data-dir>/caches/<name>`. The cache root is `<data-dir>/caches`. Contents persist between runs on the node, but are never guaranteed. An empty or partial cache dir must always be a valid starting point. Failed or killed calls keep whatever they wrote.

Cache names are global per node. Calls using the same name share contents across ops, even when they mount the cache dir at different container paths.

Concurrent calls share the directory without locking. Tools that are not safe to run concurrently, such as apt, may corrupt their cache contents.

Deleting a cache dir or the cache root is always safe when no call is using it. Later calls recreate the directories as needed.
