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
         --data "{\"host\": \"host-1\", \"service\": \"$service\", \"status\": \"ok\", \"message\": \"Service $service running\",\"nextupdatesec\": 10, \"updatekey\": \"$UPDATE_KEY\"}"

    echo -e "\n"
done

echo ""
