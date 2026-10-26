# dashGoat - Agent Instructions

## Working Rules

- Do only what the user asked. If a task reveals a problem (build error, failing test), report it and ask before changing anything.
- Always ask the user before making any change: editing files, git operations (stash, reset, checkout, commit), deleting files, or other commands that modify the working tree.
- Read-only inspection and the requested build/test commands are fine without asking.

## Build Commands

```bash
make build          # Build for current platform (runs prepare first)
make linux          # Build Linux amd64 binary
make windows        # Build Windows amd64 binary
make macintel       # Build Darwin amd64 binary
make macarm         # Build Darwin arm64 binary
make rpi            # Build Linux armv5 (Raspberry Pi)
make all            # Build all platforms
make prepare        # Copy web assets to cmd/dashgoat/web (required before build)
make clean          # Remove build artifacts and web copy
```

**Critical**: The web frontend in `/web` must be copied to `cmd/dashgoat/web` before building. Run `make prepare` or `make build` (which includes prepare).

Binary is embedded with version, commit, and build date via LDFLAGS.

## Testing

```bash
make test           # Full test suite (requires make build first)
```

Individual tests in `tests/` directory:
- `./tests/link-bin.sh` - Auto-detects OS/arch and symlinks binary to `dashgoat`
- `./tests/start-single.sh` - Starts server on port 2000 with `tests/testdashgoat.yaml`
- `./tests/stop-instances.sh` - Cleans up test processes

Test files: `ttl-test.sh`, `nextupdate-test.sh`, `metrics-test.sh`, `tags-test.sh`, `heartbeat-test.sh`, `alertmanager-test.sh`, `search-test.sh`, `dependon-test.sh`

CI runs each test suite separately with start/stop between tests.

## Architecture

- **Entry point**: `cmd/dashgoat/main.go`
- **Framework**: Echo (labstack/echo)
- **Config**: YAML file (default `dashgoat.yaml`, example: `dashgoat.example.yaml`)
- **Web assets**: `/web` directory (embedded via `//go:embed web`)
- **Prometheus**: `/metrics` endpoint (configurable via `prometheusurl`)

## Key Configuration

| Parameter | Default | Purpose |
|-----------|---------|---------|
| `updatekey` | `changeme` | API authentication for POST updates |
| `urnkey` | (empty) | Auth key for heartbeat/alertmanager GET endpoints |
| `ipport` | `:2000` | Listen address |
| `webpath` | `/` | URL prefix for all endpoints |
| `ttlbehavior` | `PromoteToOk` | Auto-resolve mode: `Remove`, `PromoteOnce`, `PromoteOneStep`, `PromoteToOk` |
| `buddy` | (empty) | Peer cluster URLs for state sync |
| `buddynsconfig` | (empty) | DNS-based buddy discovery (e.g., `dashgoat-headless-svc` in k8s) |

## API Endpoints (with webpath prefix)

- `POST /update` - Update service status
- `GET /heartbeat/:urnkey/:host/:nextupdatesec/:tags` - Watchdog heartbeat
- `POST /alertmanager/:urnkey` - Prometheus Alertmanager webhook
- `GET /status/list` - List all services
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## Environment Variables

Config values can be overridden via environment variables (same name, uppercase):
- `UPDATEKEY`, `URNKEY`, `IPPORT`, `WEBPATH`, `PROMETHEUSURL`, etc.
