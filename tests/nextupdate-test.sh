#!/bin/bash

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

# Array of test services
services=("web" "mail" "storage")

# Loop over services
for service in "${services[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"ok\", \"message\": \"Service $service running\",\"nextupdatesec\": 5, \"updatekey\": \"$UPDATE_KEY\"}"

done

echo "Waiting for nextupdatesec to expire"

sleep 10

for service in "${services[@]}"; do
    STATUS=""
    echo "Checking status for service: $service"

    STATUS=$(curl -s "$BASE_URL/status/host-1$service")

    if [ "$(echo "$STATUS" | jq -r '.status')" = "error" ]; then
        echo "Nextupdatesec is triggerd - OK"
    else
        echo "Nextupdatesec is not triggerd - ERROR"
        echo
        echo "$STATUS"
        echo
        exit 1
    fi

done