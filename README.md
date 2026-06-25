# caddy-stack

A custom [Caddy](https://caddyserver.com) build with the plugins needed to run a hardened, fully label-driven reverse proxy. One image, built reproducibly in CI, deployed on Caddy instances.

## What's compiled in

| Module | Purpose |
|---|---|
| [`caddy-docker-proxy`](https://github.com/lucaslorentz/caddy-docker-proxy) | Generates the reverse-proxy config from Docker container labels — deploy a service, label it, done. |
| [`caddy-dns/ovh`](https://github.com/caddy-dns/ovh) | Solves ACME **DNS-01** challenges via the OVH API — enables wildcard certs and certs for internal-only (tailnet) names. |
| [`caddy-crowdsec-bouncer`](https://github.com/hslatman/caddy-crowdsec-bouncer) | Enforces [CrowdSec](https://crowdsec.net) ban decisions at the proxy, before requests reach apps. |

Plugins that are **not** in the binary (they run as their own containers / update independently):

- CrowdSec **engine** (`crowdsecurity/crowdsec`) — updated via the Docker image tag.
- CrowdSec **hub content** (collections/scenarios) — updated via `cscli hub upgrade`.

## How it's built

A tiny Go module (`main.go` + `go.mod`) imports the plugins so Caddy compiles them in. `go.sum` locks the full dependency tree for reproducible builds.

```
main.go      registers the plugins via blank imports
go.mod       pins Caddy core + the three plugins
Dockerfile   multi-stage: golang builder → caddy runtime
```

The runtime base image tag is kept aligned with the Caddy core version pinned in `go.mod`.

### Build locally

```bash
docker build -t caddy-stack .
docker run --rm caddy-stack caddy list-modules | grep -E 'dns.providers.ovh|docker|crowdsec'
```

All three should be listed.

## CI / releases

GitHub Actions (`.github/workflows/build.yml`):

- **Pull requests** — build + smoke-test only (verifies all three modules are present). This is the gate that proves the image compiles *before* it reaches a server.
- **Push to `main`** — also pushes `:latest` and `:sha-<short>`.
- **Git tag `vX.Y.Z`** — pushes that semver tag for pinned rollouts.

Images are published to `ghcr.io/kamalf/caddy-stack` (public).

```bash
docker pull ghcr.io/kamalf/caddy-stack:latest
```

## Dependency updates

[Renovate](https://docs.renovatebot.com) (`renovate.json`) tracks every input:

- Caddy core, the three plugins, and the `FROM caddy` base image are grouped into one **"caddy stack"** PR.
- Minor/patch updates auto-merge once CI is green.
- The `golang` builder image is tracked automatically.

Each Renovate PR rebuilds and smoke-tests the image, so a broken upgrade never lands.

## Deployment

Used as the `caddy-docker-proxy` controller in a Compose stack. Services declare their routing via `caddy.*` labels; the controller watches the Docker socket and reconfigures Caddy automatically. DNS-01 (OVH) and the CrowdSec bouncer are configured on the controller — see the homelab project notes.
