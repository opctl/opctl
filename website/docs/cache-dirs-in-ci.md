---
title: Reuse cache dirs in CI
---

Container calls can declare [cache dirs](reference/opspec/op-directory/op/call/container/caches.md) for dependencies and compiler output. Save and restore the cache root with your CI provider's cache steps to reuse these contents between jobs. Ops must still work with an empty or partial cache dir.

The cache root is `<data-dir>/caches`. The default data dir is `$HOME/opctl` on Linux, `$HOME/Library/Application Support/opctl` on macOS, and `%LOCALAPPDATA%\opctl` on Windows. Pass `--data-dir <path>` or set `OPCTL_DATA_DIR` to relocate it. Use the same data dir for the node and the cache steps.

To check that an op works with empty cache dirs, run `opctl run --no-cache <op>`. Cache dirs start empty in scratch space shared by cache name across the run, including nested ops, loops, and parallel groups, and are removed when the run ends. It leaves the cache root untouched, so other concurrent runs can keep using their caches. Op and image caching still applies.

## GitHub Actions

This workflow assumes `opctl` is installed and the repository has an `npm-install` op at `.opspec/npm-install/op.yml` that declares an npm cache dir with a stable [cache name](reference/opspec/op-directory/op/call/container/caches.md), e.g. `caches: { /root/.npm: npm }`. Put the lockfile hash in the CI cache key, not the cache name. With a lockfile-keyed name such as `npm-$(lockHash)`, a fallback restore brings back only the old name's dir, so the op starts cold anyway. Each lockfile change also adds another full cache to the saved root, because nothing evicts old dirs. The runner also needs a container runtime supported by opctl.

```yaml
name: Install dependencies
on: [push, pull_request]

jobs:
  install:
    runs-on: ubuntu-latest
    env:
      OPCTL_DATA_DIR: ${{ runner.temp }}/opctl
    steps:
      - uses: actions/checkout@v4

      - name: Restore cache dirs
        id: cache
        uses: actions/cache/restore@v4
        with:
          path: ${{ env.OPCTL_DATA_DIR }}/caches
          key: ${{ runner.os }}-opctl-npm-${{ hashFiles('**/package-lock.json') }}
          restore-keys: |
            ${{ runner.os }}-opctl-npm-

      - name: Install dependencies
        run: opctl run npm-install

      - name: Make cache dirs readable by the runner
        if: always()
        run: |
          if [ -d "$OPCTL_DATA_DIR/caches" ]; then
            sudo chown -R "$(id -u):$(id -g)" "$OPCTL_DATA_DIR/caches"
          fi

      - name: Save cache dirs
        if: success() && steps.cache.outputs.cache-hit != 'true'
        uses: actions/cache/save@v4
        with:
          path: ${{ env.OPCTL_DATA_DIR }}/caches
          key: ${{ steps.cache.outputs.cache-primary-key }}
```

The key changes when a lockfile changes. The restore prefix lets a job start with an older cache when there is no exact match; with a stable cache name the op reuses it and the tool refreshes what changed. Old entries accumulate, so delete the CI cache occasionally if it grows too large. GitHub Actions caches are immutable, so this workflow saves only when the exact key was not restored, and only after a successful run so a partial cache isn't locked in under that key.

To limit the size of the local cache root before saving, use `opctl cache prune [NAME...]`. Give names to remove selected cache dirs, or omit them to empty the cache root. Use the same `--data-dir` or `OPCTL_DATA_DIR` as the node, and only prune when no call is using the affected dirs. It runs as root, re-executing itself with `sudo -E` when not already root (cache contents can be root-owned), and leaves `<data-dir>/nocache` untouched. Pruning locally does not change existing CI cache entries.

Containers commonly write root-owned files. When the node runs with sudo, it recursively restores each call's cache dirs to `SUDO_UID` and `SUDO_GID` after the container finishes, including failed or killed calls. A root node without those variables leaves ownership unchanged. Cleanup warnings appear in the call's stderr and do not change its outcome.

An unprivileged node cannot change ownership of files written by another user. Run `sudo chown -R "$(id -u):$(id -g)" <cache root>` before saving in that case. GitHub-hosted Linux runners have passwordless sudo. The workflow runs this step even after a failed op, before the explicit save step.

Cache dirs, including `--no-cache` dirs, start with mode `0777`. After each call, the node makes them and their contents writable by all users with `a+rwX`, so images running as different users can share a cache. Files lose their setuid and setgid bits during this permission update. An unprivileged node cannot adjust entries owned by other users, such as files written by root containers.

Parent dirs created under the data dir stay `0700`, so other host users normally cannot reach caches. More open existing parents can expose them, for example with a custom `--data-dir` under `/tmp` or the workspace. Any container mounting the same cache name can read and modify it, whatever its user. Don't put secrets in cache dirs.

## Kubernetes

The Kubernetes runtime requires `/root/opctl` as the node's data dir, backed by the `opctl` PersistentVolumeClaim shared with its container pods. Cache dirs mount as `caches/<name>` subPaths of that volume. They persist across node or pod replacement only if the volume and its contents persist. An ephemeral data-dir volume loses the cache when it is removed.

The example deployment in `sdks/go/node/containerruntime/k8s/resources.yaml` clears `/root/opctl/**` before starting the node. Remove that cleanup from your deployment to keep cache dirs across node restarts, even with a persistent volume.
