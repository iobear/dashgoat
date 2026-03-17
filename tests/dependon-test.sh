#!/bin/bash

echo
echo "-- dependon test --"
echo "-- Verifies that a service with dependon is demoted when the upstream is down --"
echo

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

cleanup() {
    curl -s --request DELETE --url $BASE_URL/service/lb-1heartbeat > /dev/null
    curl -s --request DELETE --url $BASE_URL/service/lb-1app > /dev/null
    curl -s --request DELETE --url $BASE_URL/service/host-1nginx > /dev/null
}

# --- Setup: dependency host (lb-1) is up ---
echo "Setting lb-1 to ok (dependency is healthy)"
curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"lb-1\", \"service\": \"heartbeat\", \"status\": \"ok\", \"message\": \"up\", \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

# --- Dependent service in error, dependency is up: status should stay error ---
echo "Posting host-1/nginx as error with dependon: lb-1 (dependency is up)"
curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"nginx\", \"status\": \"error\", \"message\": \"failing\", \"updatekey\": \"$UPDATE_KEY\", \"dependon\": \"lb-1\"}" > /dev/null

STATUS=$(curl -s "$BASE_URL/status/host-1nginx" | jq -r '.status')
if [[ "$STATUS" != "error" ]]; then
    echo "FAIL: Expected 'error' when dependency is up, got '$STATUS'"
    cleanup
    exit 1
fi
echo "OK: Status is 'error' while dependency (lb-1) is healthy"

# --- Bring the dependency down (fully down — only service on lb-1) ---
echo "Setting lb-1 to critical (dependency fully down)"
curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"lb-1\", \"service\": \"heartbeat\", \"status\": \"critical\", \"message\": \"down\", \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

# Re-post the dependent service to trigger dependon re-evaluation
curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"nginx\", \"status\": \"error\", \"message\": \"failing\", \"updatekey\": \"$UPDATE_KEY\", \"dependon\": \"lb-1\"}" > /dev/null

STATUS=$(curl -s "$BASE_URL/status/host-1nginx" | jq -r '.status')
if [[ "$STATUS" != "warning" ]]; then
    echo "FAIL: Expected 'warning' when dependency is fully down, got '$STATUS'"
    cleanup
    exit 1
fi
echo "OK: Status demoted to 'warning' when dependency (lb-1) is fully down"

MSG=$(curl -s "$BASE_URL/status/host-1nginx" | jq -r '.message')
if [[ "$MSG" != *"lb-1"* ]]; then
    echo "FAIL: Expected message to reference dependency 'lb-1', got: '$MSG'"
    cleanup
    exit 1
fi
echo "OK: Message references the failing dependency: $MSG"

# --- Partially down: lb-1 has one ok service and one down service ---
echo "Adding a healthy service on lb-1 (partially down scenario)"
curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"lb-1\", \"service\": \"app\", \"status\": \"ok\", \"message\": \"still up\", \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"nginx\", \"status\": \"error\", \"message\": \"failing\", \"updatekey\": \"$UPDATE_KEY\", \"dependon\": \"lb-1\"}" > /dev/null

STATUS=$(curl -s "$BASE_URL/status/host-1nginx" | jq -r '.status')
if [[ "$STATUS" != "info" ]]; then
    echo "FAIL: Expected 'info' when dependency is partially down, got '$STATUS'"
    cleanup
    exit 1
fi
echo "OK: Status demoted to 'info' when dependency (lb-1) is partially down"

# --- Recovery: dependency back to fully healthy ---
echo "Restoring lb-1 to ok (dependency recovered)"
curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"lb-1\", \"service\": \"heartbeat\", \"status\": \"ok\", \"message\": \"recovered\", \"updatekey\": \"$UPDATE_KEY\"}" > /dev/null

curl -s -X POST "$BASE_URL/update" \
     -H "Content-Type: $CONTENT_TYPE" \
     --data "{\"host\": \"host-1\", \"service\": \"nginx\", \"status\": \"error\", \"message\": \"failing\", \"updatekey\": \"$UPDATE_KEY\", \"dependon\": \"lb-1\"}" > /dev/null

STATUS=$(curl -s "$BASE_URL/status/host-1nginx" | jq -r '.status')
if [[ "$STATUS" != "error" ]]; then
    echo "FAIL: Expected 'error' after dependency recovers, got '$STATUS'"
    cleanup
    exit 1
fi
echo "OK: Status restored to 'error' after dependency (lb-1) recovered"

echo "cleaning up data"
cleanup
