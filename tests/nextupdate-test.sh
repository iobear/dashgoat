#!/bin/bash

echo
echo "-- nextupdate test --"
echo

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

for service in "${services[@]}"; do

    echo "Checking probe and change is the same ${service}"
    PROBE=$(curl -s "$BASE_URL/status/host-1${service}" | jq '.probe')
    CHANGE=$(curl -s "$BASE_URL/status/host-1${service}" | jq '.change')

    if [[ $PROBE -ne $CHANGE ]]; then
        echo $PROBE $CHANGE
        echo "Unexpected API response for .change and .probe ${service}"
        exit 1
    else
        echo "OK"
    fi
done

sleep 1

# Loop over services
for service in "${services[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"ok\", \"message\": \"Service $service running\",\"nextupdatesec\": 5, \"updatekey\": \"$UPDATE_KEY\"}"

done

for service in "${services[@]}"; do

    echo "Checking probe and change is not the same ${service}"
    PROBE=$(curl -s "$BASE_URL/status/host-1${service}" | jq '.probe')
    CHANGE=$(curl -s "$BASE_URL/status/host-1${service}" | jq '.change')

    if [[ $PROBE -gt $CHANGE ]]; then
        echo "OK"
    else
        echo $PROBE $CHANGE
        echo "Unexpected API response for .change and .probe ${service}"
        exit 1
    fi
done


echo "Waiting for nextupdatesec to expire"

sleep 11

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

echo "cleaning up data"

for service in "${services[@]}"; do
    curl -s --request DELETE --url $BASE_URL/service/host-1${service}
done