# containerd container runtime

This backend implements opctl's `ContainerRuntime` by driving **containerd**
through the [`nerdctl`](https://github.com/containerd/nerdctl) CLI. Select it with:

```sh
opctl --container-runtime containerd ...
```

The default runtime remains `docker`, so enabling this backend is opt-in and
leaves existing dockerd-based nodes/jobs untouched — the `--container-runtime`
flag value *is* the feature flag.

## Why containerd instead of dockerd

The Docker daemon's `registry-mirrors` only mirrors **Docker Hub**. It cannot
transparently redirect pulls for other registries (e.g. `quay.io`). containerd
supports **per-registry mirrors** via `hosts.toml`, which lets you redirect both
Docker Hub *and* Quay pulls through a single pull-through cache (PTC) without
rewriting any image refs in `op.yml`.

## How it differs from the docker backend

| Concern | docker backend | containerd backend |
|---|---|---|
| Transport | Docker Engine API over a socket | shells out to `nerdctl` |
| Registry mirrors | dockerd `registry-mirrors` (Hub only) | `hosts.toml` per registry (any registry) |
| Registry auth | per-op `pullCreds`, else dockerd config | containerd/nerdctl host config (`hosts.toml` + docker `config.json` cred helpers) |
| stdout/stderr | TTY | separate streams (no TTY) |

> **Auth note:** this backend does **not** use per-op `image.pullCreds`. Registry
> credentials are resolved by containerd/nerdctl host configuration (see
> [Registry auth](#registry-auth-ecr-pull-through-cache) below). This is by design:
> the whole point is that auth and mirroring are runner-level config, not per-op.
>
> **GPU note:** GPU passthrough auto-detection (a docker-backend feature) is not
> implemented here.

## Requirements

- `containerd` running on the host.
- `nerdctl` on `PATH` (or point opctl at it with `$OPCTL_NERDCTL=/path/to/nerdctl`).
- The CNI plugins nerdctl needs for bridge networking + name resolution.
- opctl running as root (rootful containerd), same as the docker backend.

opctl fails fast at node start if `nerdctl` can't be found.

## Registry mirrors (`hosts.toml`)

Point containerd at a `certs.d` directory and add one `hosts.toml` per upstream
registry. With our PTC in the artifact-archive prod account
(`012473847565`, `us-west-2`) and PTC rules `docker-hub → registry-1.docker.io`
and `quay → quay.io`:

`/etc/containerd/certs.d/docker.io/hosts.toml`:

```toml
server = "https://registry-1.docker.io"

[host."https://012473847565.dkr.ecr.us-west-2.amazonaws.com/v2/docker-hub"]
  capabilities = ["pull", "resolve"]
  # ECR's PTC path is /v2/docker-hub/<repo>; override_path stops containerd from
  # appending the upstream repo path to a path that already encodes it.
  override_path = true
```

`/etc/containerd/certs.d/quay.io/hosts.toml`:

```toml
server = "https://quay.io"

[host."https://012473847565.dkr.ecr.us-west-2.amazonaws.com/v2/quay"]
  capabilities = ["pull", "resolve"]
  override_path = true
```

And ensure containerd's config references the directory
(`/etc/containerd/config.toml`):

```toml
[plugins."io.containerd.grpc.v1.cri".registry]
  config_path = "/etc/containerd/certs.d"
```

`nerdctl` also reads `/etc/containerd/certs.d` directly. This `override_path = true`
+ `/v2/<rule>` pattern is the configuration verified to work in
[containerd discussion #11385](https://github.com/containerd/containerd/discussions/11385).

Ready-to-copy versions of these files live in [`examples/`](examples).

## Registry auth (ECR pull-through cache)

The PTC registry is private ECR, so pulls must authenticate. **Preferred** is the
[`docker-credential-ecr-login`](https://github.com/awslabs/amazon-ecr-credential-helper)
helper, keyed to the ECR registry host in the docker config that nerdctl reads
(`$DOCKER_CONFIG/config.json`, default `~/.docker/config.json`):

```json
{
  "credHelpers": {
    "012473847565.dkr.ecr.us-west-2.amazonaws.com": "ecr-login"
  }
}
```

**Auth-refresh mechanism:** ECR authorization tokens expire after 12 hours.
`docker-credential-ecr-login` is invoked by nerdctl *on every pull* for the ECR
host; it calls `ecr:GetAuthorizationToken` using the runner's AWS credential
chain (its service role) and returns a fresh token (caching it until shortly
before expiry). There is **no token to rotate manually** and nothing to schedule
— keep the helper on `PATH` and the role's credentials available.

Because pulls are redirected to the mirror host, nerdctl must attach these
credentials to the **mirror host** (the ECR registry), not the original
`docker.io`/`quay.io` host. Recent containerd/nerdctl do this — the resolver
authenticates against the host it actually connects to. Verify on the target
runner once:

```sh
nerdctl pull --debug-full docker.io/library/alpine:3.20 2>&1 \
  | grep -i 'authorized\|012473847565.dkr.ecr'
```

**Fallback** (only if your nerdctl version does *not* apply host creds on the
mirror path): put a static `Authorization` header in `hosts.toml` and refresh it
out of band. ECR Basic auth is `base64("AWS:<token>")`:

```toml
[host."https://012473847565.dkr.ecr.us-west-2.amazonaws.com/v2/docker-hub".header]
  Authorization = "Basic QVdTOjxFQ1IgdG9rZW4+"
```

Refresh it with a systemd timer (every ~6h, well inside the 12h expiry):

```sh
TOKEN="$(aws ecr get-login-password --region us-west-2)"
AUTH="Basic $(printf 'AWS:%s' "$TOKEN" | base64 -w0)"
# rewrite the Authorization line in each hosts.toml, then no daemon restart needed
```

## Proxy bypass (`NO_PROXY`)

Our runners have no direct internet route; egress goes through a mandatory
forward proxy. The ECR/S3 endpoints used by the PTC are reachable via interface
VPC endpoints and **must bypass** the proxy, or pulls will fail. containerd and
nerdctl honor `HTTPS_PROXY`/`NO_PROXY` from **their own** process environment,
so set `NO_PROXY` on the **containerd service** (a systemd drop-in), not just in
the build shell:

```ini
# /etc/systemd/system/containerd.service.d/http-proxy.conf
[Service]
Environment="HTTPS_PROXY=http://forward-proxy:3128"
Environment="NO_PROXY=*.dkr.ecr.us-west-2.amazonaws.com,api.ecr.us-west-2.amazonaws.com,*.s3.us-west-2.amazonaws.com,localhost,127.0.0.1"
```

(opctl separately propagates `HTTP(S)_PROXY`/`NO_PROXY` from the node's env into
the op containers it runs — see `containerruntime.ProxyEnvVars` — so workloads
inside the containers also reach the network through the proxy.)

## IAM (runner service role)

The runner's service role needs:

- `ecr:GetAuthorizationToken` on `*`
- pull + `BatchImportUpstreamImage` on
  `arn:aws:ecr:us-west-2:012473847565:repository/{docker-hub,quay}/*`

See [`examples/iam-policy.json`](examples/iam-policy.json).

## Runner wiring (codebuild-infra-pipeline)

The opctl code change lives here; the runner-template change lives in
`Remitly/codebuild-infra-pipeline`
(`runners/default/cfn/gha-codebuild-runner.yaml`). Per the golden-pipeline
sign-off norm, that change is proposed for review rather than landed silently.
The runner's user-data / buildspec must, behind the same feature flag:

1. Install `containerd`, `nerdctl`, CNI plugins, and `docker-credential-ecr-login`.
2. Drop in the `certs.d/*/hosts.toml`, docker `config.json`, and containerd
   `config.toml` from [`examples/`](examples).
3. Add the `NO_PROXY` systemd drop-in above and `systemctl daemon-reload &&
   systemctl restart containerd`.
4. Attach the IAM policy delta to the runner service role.
5. Run opctl with `--container-runtime containerd`.

## Verifying (definition of done)

1. Run an op whose image is a Docker Hub image and one whose image is a `quay.io`
   image; both should pull and run.
2. Confirm new repositories appear under `docker-hub/…` and `quay/…` in the
   artifact-archive ECR (PTC populated them on first pull).
3. Confirm runner egress logs show **no** connections to `registry-1.docker.io`
   or `quay.io`.
