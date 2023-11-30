#!/bin/bash

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

services=("web" "mail" "storage")

for service in "${services[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"error\", \"message\": \"Service $service running\",\"ttl\": 10, \"updatekey\": \"$UPDATE_KEY\"}"

    echo -e "\n"
done

echo ""
