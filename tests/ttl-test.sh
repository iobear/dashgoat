#!/bin/bash

echo
echo "-- ttl test --"
echo

CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

services=("web" "mail" "storage")

dashgoat_bin() {
    if [ -n "$GITHUB_WORKSPACE" ]; then
        echo "$GITHUB_WORKSPACE/dashgoatbuild/dashgoat"
    else
        echo "./dashgoat"
    fi
}

start_dashgoat() {
    local port=$1
    local mode=$2
    echo "Starting dashgoat on :$port with TTLBEHAVIOR=$mode"
    TTLBEHAVIOR=$mode $(dashgoat_bin) -ipport ":$port" -configfile tests/testdashgoat.yaml > /dev/null 2>&1 &
    echo $! > PID_ttlmode

    for i in {1..10}; do
        sleep 1
        if [ "$(curl -s "http://localhost:$port/health" | jq -r '.Ready')" = "true" ]; then
            echo "dashGoat ready (mode=$mode)"
            return 0
        fi
    done
    echo "FAIL: dashGoat did not become ready (mode=$mode)"
    exit 1
}

stop_dashgoat() {
    if [ -f PID_ttlmode ]; then
        kill "$(cat PID_ttlmode)" 2>/dev/null
        rm -f PID_ttlmode
        sleep 1
    fi
}

# -----------------------------------------------------------------------
echo "=== Mode: PromoteToOk (default, port 2000) ==="

BASE_URL="http://localhost:2000"

for service in "${services[@]}"; do
    echo "Updating status for service: $service"
    curl -s -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"error\", \"message\": \"Service $service running\", \"ttl\": 5, \"updatekey\": \"$UPDATE_KEY\"}"
    if [ $? -ne 0 ]; then
        echo "Error updating status for service: $service"
        exit 1
    fi
done

echo "Waiting for TTL to expire (11s)..."
sleep 11

for service in "${services[@]}"; do
    echo "Checking status for service: $service"
    STATUS=$(curl -s "$BASE_URL/status/host-1$service")
    if [ "$(echo "$STATUS" | jq -r '.status')" = "ok" ]; then
        echo "OK: TTL expired — status promoted to ok"
    else
        echo "FAIL: TTL did not expire for $service"
        exit 1
    fi
done

echo "cleaning up data"
for service in "${services[@]}"; do
    curl -s --request DELETE --url "$BASE_URL/service/host-1${service}" > /dev/null
done

# -----------------------------------------------------------------------
PORT=2000
BASE_URL="http://localhost:$PORT"

echo
echo "=== Mode: Remove ==="
start_dashgoat $PORT "remove"

curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"web\", \"status\": \"error\", \"message\": \"test\", \"ttl\": 5, \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

echo "Waiting for TTL to expire (10s)..."
sleep 10

RESPONSE=$(curl -s "$BASE_URL/status/list")
if [[ -z "$RESPONSE" ]] || [[ "$(echo "$RESPONSE" | jq 'has("host-1web")')" == "false" ]]; then
    echo "OK: Service was removed after TTL expiry"
else
    echo "FAIL: Service still present after TTL expiry in Remove mode"
    echo "$RESPONSE"
    stop_dashgoat
    exit 1
fi

stop_dashgoat

# -----------------------------------------------------------------------
echo
echo "=== Mode: PromoteOnce ==="
start_dashgoat $PORT "promoteonce"

curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"web\", \"status\": \"critical\", \"message\": \"test\", \"ttl\": 5, \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

echo "Waiting for first TTL step (10s)..."
sleep 10

STATUS=$(curl -s "$BASE_URL/status/host-1web" | jq -r '.status')
if [[ "$STATUS" != "error" ]]; then
    echo "FAIL: Expected 'error' after one PromoteOnce step (critical → error), got '$STATUS'"
    stop_dashgoat
    exit 1
fi
echo "OK: critical → error after first TTL step"

echo "Confirming no further promotion (10s)..."
sleep 10

STATUS=$(curl -s "$BASE_URL/status/host-1web" | jq -r '.status')
if [[ "$STATUS" != "error" ]]; then
    echo "FAIL: Expected status to remain 'error' (PromoteOnce stops after one step), got '$STATUS'"
    stop_dashgoat
    exit 1
fi
echo "OK: Status stayed at 'error' — PromoteOnce does not promote further"

stop_dashgoat

# -----------------------------------------------------------------------
echo
echo "=== Mode: PromoteOneStep ==="
start_dashgoat $PORT "promoteonestep"

curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"web\", \"status\": \"critical\", \"message\": \"test\", \"ttl\": 5, \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

echo "Waiting for first TTL step (10s)..."
sleep 10

STATUS=$(curl -s "$BASE_URL/status/host-1web" | jq -r '.status')
if [[ "$STATUS" != "error" ]]; then
    echo "FAIL: Expected 'error' after first PromoteOneStep (critical → error), got '$STATUS'"
    stop_dashgoat
    exit 1
fi
echo "OK: First step: critical → error"

echo "Waiting for second TTL step (10s)..."
sleep 10

STATUS=$(curl -s "$BASE_URL/status/host-1web" | jq -r '.status')
if [[ "$STATUS" != "warning" ]]; then
    echo "FAIL: Expected 'warning' after second PromoteOneStep (error → warning), got '$STATUS'"
    stop_dashgoat
    exit 1
fi
echo "OK: Second step: error → warning"

curl -s --request DELETE --url "$BASE_URL/service/host-1web" > /dev/null
stop_dashgoat

echo
echo "All TTL tests passed"