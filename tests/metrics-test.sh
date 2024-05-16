#!/bin/bash

echo
echo "-- metrics test --"
echo

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

services=("web" "mail" "storage")

for service in "${services[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"error\", \"message\": \"Service $service running\", \"updatekey\": \"$UPDATE_KEY\"}"

done

for service in "${services[@]}"; do
    STATUS=""
    echo "Checking status for service: $service"

    STATUS=$(curl -s "$BASE_URL/metrics")

    if [ "$(echo "$STATUS" | grep -c $service)" = "1" ]; then
        echo "$service found at /metrics"
    else
        echo "$service not found at /metrics - ERROR"
        exit 1
    fi

done

echo "Test delete service in prometheus"

curl --request DELETE \
  --url $BASE_URL/service/host-1web

STATUS=$(curl -s "$BASE_URL/metrics")

if [ "$(echo "$STATUS" | grep -c host-1web)" = "0" ]; then
    echo "$service not found at /metrics - OK"
else
    echo "$service found at /metrics - ERROR"
    exit 1
fi

echo "cleaning up data"

curl -s --request DELETE --url $BASE_URL/service/host-1mail
curl -s --request DELETE --url $BASE_URL/service/host-1storage
