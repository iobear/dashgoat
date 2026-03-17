# dashGoat API

dashGoat exposes a simple HTTP/JSON API for submitting and querying service status, heartbeats, and alert integrations.

## Base URL

```
http://localhost:2000
```

The base path can be changed via the `webpath` config option.

---

## Authentication

Two key types are used depending on the endpoint:

| Key | Config / Flag | Env var | Used by |
|-----|--------------|---------|---------|
| `updatekey` | `-updatekey` / `updatekey:` | `UPDATEKEY` | `POST /update` (in request body) |
| `urnkey` | `-urnkey` / `urnkey:` | `URNKEY` | `GET /heartbeat/…`, `POST /alertmanager/…` (in URL path) |

Read (GET) and delete endpoints require no authentication.

---

## Endpoints

### Status updates

#### `POST /update`

Submit or update the status of a host/service pair.

**Request body** (`application/json`):

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `updatekey` | string | ✅ | API write key |
| `host` | string | ✅ | Host identifier |
| `service` | string | ✅ | Service name |
| `status` | string | ✅ | `ok` · `info` · `warning` · `error` · `critical` |
| `message` | string | | Human-readable description |
| `severity` | string | | Override display severity |
| `nextupdatesec` | int | | Watchdog interval in seconds — warn if no update arrives within this window |
| `ttl` | int | | Auto-resolve after N seconds (see TTL modes in README) |
| `tags` | string[] | | Labels for grouping / dependency matching |
| `probe` | int | | Last-seen Unix timestamp (set by probe, not dashGoat) |
| `from` | string[] | | Source identifiers |
| `dependon` | string | | Host or tag this service depends on — reduces severity when upstream is down |

**Response:** `201 Created` — the stored `ServiceState` object.

**Example:**

```bash
curl -X POST http://localhost:2000/update \
  -H 'Content-Type: application/json' \
  -d '{
    "host": "host-1",
    "service": "nginx",
    "status": "ok",
    "message": "All good",
    "updatekey": "my-precious!",
    "nextupdatesec": 60,
    "tags": ["production", "uk"]
  }'
```

---

### Heartbeat (watchdog via GET)

#### `GET /heartbeat/:urnkey/:host/:nextupdatesec/:tags`

Lightweight keep-alive endpoint, useful when only HTTP GET is available (e.g. routers, embedded devices). Sets the service `heartbeat` on the given host to `ok`.

| Param | Description |
|-------|-------------|
| `urnkey` | URN authentication key |
| `host` | Host identifier |
| `nextupdatesec` | Watchdog interval in seconds |
| `tags` | Comma-separated tags (e.g. `router,home`) |

**Response:** `200 OK` — the stored `ServiceState` object.

**Example:**

```bash
curl http://localhost:2000/heartbeat/my-urnkey/router01/30/router,home
```

---

### Read status

#### `GET /status/list`

Returns all tracked service states as a JSON object keyed by `{host}{service}`.

**Response:** `200 OK` · `204 No Content` (when empty)

```json
{
  "demo2buddy": {
    "service": "buddy",
    "host": "demo2",
    "status": "ok",
    "message": "buddy up",
    "severity": "info",
    "nextupdatesec": 0,
    "tags": null,
    "probe": 1718304779,
    "change": 1718304779,
    "from": ["demo"],
    "ack": "",
    "ttl": 0,
    "dependon": "",
    "UpdateKey": "valid"
  }
}
```

---

#### `GET /status/:id`

Returns a single service state by its ID (`{host}{service}`).

**Response:** `200 OK`

---

#### `GET /status/listmso`

Returns a simplified status list in MSO/HPO format.

**Response:** `200 OK`

```json
{
  "host-1nginx": {
    "status": "ok",
    "message": "[ok] nginx host-1 - All good"
  }
}
```

---

#### `GET /status/listsearch/:search_string`

Search service states by host or service name. Allowed characters: `a-z A-Z 0-9 + , - _`

**Response:** `200 OK` · `204 No Content` (no matches)

**Example:**

```bash
curl http://localhost:2000/status/listsearch/nginx
```

---

### List unique values

#### `GET /list/:field`

Returns a sorted array of unique values for a given field across all service states.

Common values for `field`: `host`, `service`, `status`, `tags`

**Response:** `200 OK`

**Examples:**

```bash
curl http://localhost:2000/list/host
# ["demo2","gateway-1","mailserver-1","nas-1"]

curl http://localhost:2000/list/status
# ["error","ok"]
```

---

### Delete

#### `DELETE /service/:id`

Removes a service state. The ID is `{host}{service}` (spaces converted to hyphens).

**Response:** `204 No Content` · `404 Not Found`

**Example:**

```bash
curl -X DELETE http://localhost:2000/service/host-1nginx
```

---

### Alertmanager integration

#### `POST /alertmanager/:urnkey`

Accepts a Prometheus Alertmanager webhook payload and converts it to dashGoat service states. Requires `urnkey` to be configured.

**Request body:** Standard [Alertmanager webhook payload](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config) (version 4).

**Response:** `201 Created`

**Alertmanager config snippet:**

```yaml
receivers:
  - name: dashgoat
    webhook_configs:
      - send_resolved: true
        url: https://<dashgoat-host>/alertmanager/<urnkey>
route:
  receiver: dashgoat
  repeat_interval: 5m
  matchers:
    - alertname != Watchdog
```

---

### Observability

#### `GET /health`

Returns `200 OK` when dashGoat is ready; `503 Service Unavailable` otherwise.

**Response body:**

```json
{
  "Hostnames": ["demo_65.109.171.87", "demo"],
  "DashName": "demo",
  "Ready": true,
  "UpAt": "2024-06-13T18:52:47.039958109Z",
  "UpAtEpoch": 1718304767,
  "DashGoatVersion": "v1.7.8",
  "GoVersion": "go1.22.4",
  "BuildDate": "2024-06-09",
  "MetricsHistory": true,
  "Prometheus": true,
  "Commit": "dee8dac"
}
```

---

#### `GET /metrics`

Exposes Prometheus metrics (can be disabled with `disablemetrics: true` in config).

---

#### `GET /metricshistory/:serviceid/:hours`

Fetches historical status data for a service from a connected Prometheus instance.

| Param | Description |
|-------|-------------|
| `serviceid` | Service identifier (`{host}{service}`) |
| `hours` | Number of hours to look back |

**Requirements:** `prometheusurl` must be configured.

**Response:** `200 OK` — array of `ServiceStatus` with timestamps.

---

## Status values

| Value | Meaning |
|-------|---------|
| `ok` | Service is healthy |
| `info` | Informational, no action needed |
| `warning` | Degraded, attention recommended |
| `error` | Service is failing |
| `critical` | Service is down / critical failure |

---

## Full request reference

A Postman / Insomnia collection is available at `doc/dashGoat.postman_collection.json`.
