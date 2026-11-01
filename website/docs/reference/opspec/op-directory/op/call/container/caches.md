---
title: Caches [object]
---

An object mapping absolute container paths to cache names and sharing options.

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

An explicit value can be a string literal or [string expression](../../../../types/string.md#initialization), such as `npm-$(lockHash)`, evaluated in the call's scope. Every cache name, including a derived name, must be non-empty and may contain only ASCII letters, digits, `.`, `_`, and `-`. The names `.` and `..` are invalid. A container path cannot appear in both `caches` and `dirs`, `files`, or a `sockets` entry that bind-mounts a Unix socket.

Opspec has no hash function. Compute a lockfile hash in a prior call, then use it in the cache name when contents must not be shared across lockfile versions. This example writes the hash to a file bound to `lockHash`, which the next call reads through string interpolation. The example assumes `op.yml` lives at `.opspec/npm-install/` in the project. Remove the trailing newline so it does not become part of the name.

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

Each distinct name is a separate cache dir that the node never evicts, so a hash-keyed name leaves a full cache behind every time the hash changes. Prefer a stable name, such as `npm`, for caches that tolerate mixed contents, like npm's content-addressed cache. Reserve hash-keyed names for caches that must not be shared across versions, and prune old dirs with `opctl cache prune NAME...` when no call is using them.

The node stores each cache dir at `<data-dir>/caches/<name>`. The cache root is `<data-dir>/caches`. Contents persist between runs on the node, but are never guaranteed. An empty or partial cache dir must always be a valid starting point. Failed or killed calls keep whatever they wrote.

With the embedded runtime, Lima shares the host's home directory with the VM at the same path and with write access. The default data dir is under that share, so cache dirs and `--no-cache` dirs are visible to Docker in the VM. Keep a custom `--data-dir` under the home directory too. Docker in the VM can't see other host paths, apart from Lima's own `/tmp/lima` share.

On Kubernetes, the data dir must be `/root/opctl` on the `opctl` PersistentVolumeClaim; cache and `--no-cache` dirs mount as subPaths of it. See [Kubernetes CI guidance](../../../../../../cache-dirs-in-ci.md#kubernetes).

Use `opctl run --no-cache <op>` for a cold run. The run starts every cache dir empty in scratch space under `<data-dir>/nocache/<run id>`, shared by cache name across the run's calls, including calls in nested ops, loops, and parallel groups, and removed when the run ends. A killed run's leftovers are swept when the node next starts. A node run without sudo can't remove root-owned files that containers wrote, so those stay until you remove them with the node stopped, e.g. `sudo rm -rf <data-dir>/nocache`. The run leaves the cache root untouched. Other concurrent runs on the node keep using their caches. Op and image caching still applies.

Cache dirs, including `--no-cache` dirs, start with mode `0777`. After each call, the node makes them and their contents writable by all users with `a+rwX`, so images running as different users can share a cache. Files lose their setuid and setgid bits during this permission update. An unprivileged node cannot adjust entries owned by other users, such as files written by root containers.

Parent dirs created under the data dir stay `0700`, so other host users normally cannot reach caches. More open existing parents can expose them, for example with a custom `--data-dir` under `/tmp` or the workspace. Any container mounting the same cache name can read and modify it, whatever its user. Don't put secrets in cache dirs.

After every container call, including failed or killed calls, a node run with sudo recursively restores that call's cache dir ownership to `SUDO_UID` and `SUDO_GID`. A node run as root without sudo leaves ownership unchanged. It does not follow symlinks. Cleanup errors produce warnings in the call's stderr without changing its result. An unprivileged node also leaves ownership unchanged, so container-written files may remain root-owned. See [reuse cache dirs in CI](../../../../../../cache-dirs-in-ci.md) for save and restore steps, including the manual ownership fix.

Cache names are global per node. Calls using the same name share contents across ops, even when they mount the cache dir at different container paths. On case-insensitive filesystems (the macOS and Windows defaults), names that differ only by case share one cache dir.

Cache dirs use `sharing: shared` by default. For tools that cannot safely share a cache concurrently, such as apt, use the object form with `sharing: locked`:

```yaml
run:
  container:
    image: { ref: debian:stable }
    caches:
      /var/cache/apt: { name: apt, sharing: locked }
      /var/lib/apt: { name: apt-lists, sharing: locked }
    cmd: [apt-get, update]
```

The optional `name` accepts the same string literals and expressions as the string form. Omitting it or setting it to `null` derives the name from the image repository and container path. `sharing` accepts only `shared` and `locked`; omitted `sharing` means `shared`. Other object keys are rejected.

Calls with `sharing: locked` wait for other locked calls using the same cache dir on the node. The lock lasts through the whole container run and the node's permission and ownership cleanup, including after failure or kill. Shared calls neither wait for nor block locked calls, so every call using a concurrency-sensitive cache should opt into locking. A waiting call reports the cache name in stderr and can be killed. Calls locking several caches acquire them in a consistent order. If one call mounts the same cache more than once with different `sharing` values, it locks. On case-insensitive filesystems, locking follows the exact spelling of the name, so give every locked user the same spelling.

Locking also applies within a `--no-cache` run. Its scratch dirs have separate locks from the warm cache and other cold runs.

Deleting a cache dir or the cache root is always safe when no call is using it. Later calls recreate the directories as needed.

Use `opctl cache prune [NAME...]` to remove named cache dirs, or omit names to empty the whole cache root. It uses `--data-dir` or `OPCTL_DATA_DIR`, runs as root (re-executing itself with `sudo -E` when not already root, since cache contents can be root-owned), and ignores missing names. All names must pass the cache name rules above before anything is deleted. Only prune when no call is using the affected dirs; the command does not check for active calls. It leaves `<data-dir>/nocache` untouched.
