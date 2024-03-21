#!/bin/bash

echo
echo "-- heartbeat test --"

BASE_URL="http://localhost:2000"
HEARTBEAT_KEY="TmwW8kOO1Cmks54TA"

hosts=("web-1" "mail-1" "storage-1")

for host in "${hosts[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/heartbeat/$HEARTBEAT_KEY/$host/5/help" \

done


echo "Waiting for nextupdatesec to expire"

sleep 11

for host in "${hosts[@]}"; do
    STATUS=""
    echo "Checking status for host: $host"

    STATUS=$(curl -s "$BASE_URL/status/${host}heartbeat")

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