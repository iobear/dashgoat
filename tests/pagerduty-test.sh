#!/bin/bash

echo
echo "-- pagerduty test --"

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

# Define services with tags
services_with_tags=(
    '"service": "nginx", "tags": ["web", "production"]'
    '"service": "database", "tags": ["db", "production"]'
    '"service": "cache", "tags": ["cache", "development"]'
)


for service_with_tags in "${services_with_tags[@]}"; do

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", $service_with_tags, \"status\": \"error\", \"message\": \"Slow service\",\"updatekey\": \"$UPDATE_KEY\"}"

done
