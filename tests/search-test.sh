#!/bin/bash

echo
echo "-- search test --"
echo

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

services=("web" "mail" "storage")

for service in "${services[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"ok\", \"message\": \"Service $service running\", \"updatekey\": \"$UPDATE_KEY\"}"
    if [ $? -ne 0 ]; then
        echo "Error updating status for service: $service"
        exit 1
    fi

done

echo "Searching for service mail"

SERVICE=$(curl -s "$BASE_URL/status/listsearch/mail")


if [ "$(echo "$SERVICE" | jq -r '.[] | .service')" = "mail" ]; then
    echo "Service mail found - OK"
else
    echo "Service mail not found - ERROR"
    exit 1
fi

echo "Searching for service stora expecting to match storage"

SERVICE=$(curl -s "$BASE_URL/status/listsearch/stora")
if [ "$(echo "$SERVICE" | jq -r '.[] | .service')" = "storage" ]; then
    echo "Service storage found - OK"
else
    echo "Service storage not found - ERROR"
    exit 1
fi

echo "cleaning up data"

for service in "${services[@]}"; do
    curl -s --request DELETE --url $BASE_URL/service/host-1${service}
done