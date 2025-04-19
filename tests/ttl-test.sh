#!/bin/bash

echo
echo "-- ttl test --"
echo

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

services=("web" "mail" "storage")

for service in "${services[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"error\", \"message\": \"Service $service running\",\"ttl\": 5, \"updatekey\": \"$UPDATE_KEY\"}"
    if [ $? -ne 0 ]; then
        echo "Error updating status for service: $service"
        exit 1
    fi

done

echo "Waiting for TTL to expire"

sleep 11

for service in "${services[@]}"; do
    STATUS=""
    echo "Checking status for service: $service"

    STATUS=$(curl -s "$BASE_URL/status/host-1$service")

    if [ "$(echo "$STATUS" | jq -r '.status')" = "ok" ]; then
        echo "TTL is expired - OK"
    else
        echo "TTL is not expired - ERROR"
        exit 1
    fi

done

echo "cleaning up data"

for service in "${services[@]}"; do
    curl -s --request DELETE --url $BASE_URL/service/host-1${service}
done