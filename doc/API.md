# dashGoat API Documentation

dashGoat API documentation! Below you will find a list of available endpoints, including required headers, request bodies, and examples.

## Base URL

The base URL for all API requests is `http://localhost:2000`.

## Authentication

All requests to the dashGoat API require an `updatekey` for authentication. Include your `updatekey` in the request body as shown in the examples below.

## Content-Type

All requests should include the `Content-Type: application/json` header unless otherwise specified.

## API Endpoints

### Update Status

- **POST** `/update`
  - **Description**: Updates the status of a given host and service.
  - **Headers**:
    - `Content-Type: application/json`
  - **Body**:
    ```json
    {
      "host": "host-1",
      "service": "HTTP",
      "status": "ok", // <ok/info/warning/error/critical>
      "message": "Hello World",
      "severity": "info", // (optional)
      "nextupdatesec": 601, // expect update within 601 seconds (optional)
      "ttl": 600, // If no new update, assume system ok after 600 seconds (optional)
      "tags": ["production","uk","customer"], // (optional)
      "probe": 1610839637, // last seen unix timestamp (optional)
      "from": ["uptimeprobe-uk"], // (optional)
      "dependon": "loadbalancer-uk-1", // My service depends on loadbalancer-uk-1 (optional)
      "updatekey": "your_updatekey_here"
    }
    ```
  - **Notes**: Replace `your_updatekey_here` with your actual `dashgoat_updatekey`.
  <br/>(*) Is optional

### Health Check

- **GET** `/health`
  - **Description**: Checks the health of the dashGoat service.
  - **Headers**:
    - `Content-Type: application/json`

### Status List

- **GET** `/status/list`
  - **Description**: Retrieves a list of all status updates.
  - **Headers**:
    - `Accept: application/json`
```json
{
"demo2buddy":
  {"service":"buddy","host":"demo2","status":"ok","message":"buddy up","severity":"info","nextupdatesec":0,"tags":null,"probe":1718304779,"change":1718304779,"from":["demo"],"ack":"","ttl":0,"dependon":"","UpdateKey":"valid"},
"gateway-1heartbeat":
  {"service":"heartbeat","host":"gateway-1","status":"ok","message":"","severity":"info","nextupdatesec":66,"tags":["router","demo"],"probe":1718306581,"change":1718304781,"from":["heartbeat"],"ack":"","ttl":0,"dependon":"","UpdateKey":"valid"}
}
```

### List Known Hosts

- **GET** `/list/host`
  - **Description**: Retrieves a list of all known hosts.
  - **Headers**:
    - `Accept: application/json`
```json
[
 "demo2",
 "gateway-1",
 "mailserver-1",
 "nas-1"
]
```

### List Known Services

- **GET** `/list/status`
  - **Description**: Retrieves a list of all current status'
  - **Headers**:
    - `Accept: application/json`

```json
[
 "ok",
 "error"
]
```

### Delete Service Status

- **DELETE** `/service/{host}{service}`
  - **Description**: Deletes the status of a specific service for a host.

### Metrics

- **Metrics** `/metrics`
  - **Description**: Shows Prometheus metrics

### Health

- **Health** `/health`
  - **Description**: Returns HTTP 200 when ready and healthy, A JSON object is also returned

  ```json
  {
   "Hostnames":["demo_65.109.171.87","demo"],
   "DashName":"demo",
   "Ready":true,
   "UpAt":"2024-06-13T18:52:47.039958109Z",
   "UpAtEpoch":1718304767,
   "DashGoatVersion":"v1.7.8",
   "GoVersion":"go1.22.4","
   BuildDate":"2024-06-09",
   "MetricsHistory":true,
   "Prometheus":true,
   "Commit":"dee8dac"
  }
  ```

## Getting Started

To get started with the dashGoat API, ensure you have the correct `dashgoat_updatekey` and make a request to one of the endpoints listed above using the specified method, headers, and body format.
