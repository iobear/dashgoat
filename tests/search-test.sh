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

echo "Searching by host name: host-1 should match all three services"

SERVICE=$(curl -s "$BASE_URL/status/listsearch/host-1")
COUNT=$(echo "$SERVICE" | jq '[.[] | .host] | length')
if [ "$COUNT" -eq "${#services[@]}" ]; then
    echo "All ${#services[@]} services found via host search - OK"
else
    echo "Expected ${#services[@]} results searching by host, got $COUNT - ERROR"
    exit 1
fi

echo "Searching by partial host name: host should match all services on host-1"

SERVICE=$(curl -s "$BASE_URL/status/listsearch/host")
COUNT=$(echo "$SERVICE" | jq '[.[] | .host] | length')
if [ "$COUNT" -ge "${#services[@]}" ]; then
    echo "Services matched via partial host name - OK"
else
    echo "Expected at least ${#services[@]} results for partial host search, got $COUNT - ERROR"
    exit 1
fi

echo "Searching for non-existent term should return no results"

RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/status/listsearch/zzz-no-match")
if [ "$RESPONSE" = "204" ] || [ "$RESPONSE" = "200" ]; then
    echo "Non-matching search returned $RESPONSE - OK"
else
    echo "Expected 204 or 200 for non-matching search, got $RESPONSE - ERROR"
    exit 1
fi

echo "cleaning up data"

for service in "${services[@]}"; do
    curl -s --request DELETE --url $BASE_URL/service/host-1${service}
done