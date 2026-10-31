---
title: Caches [object]
---

An object mapping absolute container paths to cache names.

```yaml
run:
  container:
    image: { ref: golang:1.23 }
    caches:
      /root/.cache/go-build:
      /go/pkg/mod: gomod
```

Leave a value empty, or use `null`, to derive the cache name from the interpreted image reference and container path. The node lowercases the image reference, drops its tag and digest, and uses its familiar repository name, so `docker.io/library/node` becomes `node`. It appends the cleaned container path, replaces characters outside `[A-Za-z0-9._-]` with `-`, collapses runs of `-`, and trims leading and trailing `-`.

| Image reference | Container path | Derived cache name |
| --- | --- | --- |
| `node:22` | `/root/.npm` | `node-root-.npm` |
| `golang:1.23` | `/go/pkg/mod/` | `golang-go-pkg-mod` |
| `ghcr.io/org/tool:1` | `/cache` | `ghcr.io-org-tool-cache` |
| `localhost:5000/x/y:tag` | `/c` | `localhost-5000-x-y-c` |

Tags and digests of the same repository share a cache at the same container path. Different image repositories usually get different cache names, but sanitizing can make distinct repository and path pairs share one; give an explicit name when that matters. Give images the same explicit name to share a cache across images. Images loaded from an OCI image directory and bare image IDs have no repository reference and require an explicit cache name.

An explicit value can be a string literal or [string expression](../../../../types/string.md#initialization), such as `npm-$(lockHash)`, evaluated in the call's scope. Every cache name, including a derived name, must be non-empty and may contain only ASCII letters, digits, `.`, `_`, and `-`. The names `.` and `..` are invalid. A container path cannot appear in both `dirs` and `caches`.

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

With the embedded runtime, Lima shares the host's home directory with the VM at the same path and with write access. The default data dir is under that share, so cache dirs and `--no-cache` dirs are visible to Docker in the VM. Keep a custom `--data-dir` under the home directory too. Paths outside the home directory aren't visible to Docker in the VM.

On Kubernetes, the data dir must be `/root/opctl` on the `opctl` PersistentVolumeClaim; cache and `--no-cache` dirs mount as subPaths of it. See [Kubernetes CI guidance](../../../../../../cache-dirs-in-ci.md#kubernetes).

Use `opctl run --no-cache <op>` for a cold run. The run starts every cache dir empty in scratch space under `<data-dir>/nocache/<run id>`, shared by cache name across the run's calls, including calls in nested ops, loops, and parallel groups, and removed when the run ends. The run leaves the cache root untouched. Other concurrent runs on the node keep using their caches. Op and image caching still applies.

Cache dirs are created with mode `0700`, owned by the user and group who ran sudo, or the node's effective user and group otherwise. An image running as a different non-root `USER` may be unable to write them, as with existing scratch dirs.

After every container call, including failed or killed calls, a node run with sudo recursively restores that call's cache dir ownership to `SUDO_UID` and `SUDO_GID`. A node run as root without sudo leaves ownership unchanged. It does not follow symlinks. Ownership errors produce warnings in the call's stderr without changing its result. An unprivileged node also leaves ownership unchanged, so container-written files may remain root-owned. See [reuse cache dirs in CI](../../../../../../cache-dirs-in-ci.md) for save and restore steps, including the manual ownership fix.

Cache names are global per node. Calls using the same name share contents across ops, even when they mount the cache dir at different container paths.

Concurrent calls share the directory without locking. Tools that are not safe to run concurrently, such as apt, may corrupt their cache contents.

Deleting a cache dir or the cache root is always safe when no call is using it. Later calls recreate the directories as needed.
