#!/bin/bash

echo
echo "-- tags test --"

BASE_URL="http://localhost:2000"
CONTENT_TYPE="application/json"
UPDATE_KEY="changeme"

# Define services with tags
services_with_tags=(
    '"service": "nginx", "tags": ["web", "production"]'
    '"service": "database", "tags": ["db", "production"]'
    '"service": "cache", "tags": ["cache", "development"]'
    '"service": "storage", "tags": ["storage", "development"]'
)

for service_with_tags in "${services_with_tags[@]}"; do

    curl -X POST "$BASE_URL/update" \
         -H "Content-Type: $CONTENT_TYPE" \
         --data "{\"host\": \"host-1\", $service_with_tags, \"status\": \"ok\", \"message\": \"Service is up\",\"updatekey\": \"$UPDATE_KEY\"}"

done

for service_with_tags in "${services_with_tags[@]}"; do
    service=$(echo {"$service_with_tags"} | jq -r '.service')
    tags=$(echo {"$service_with_tags"} | jq -r '.tags[]')

#    echo "Validating tags for service: $service"

    STATUS=$(curl -s "$BASE_URL/status/host-1$service")

    for tag in $tags; do
        if echo "$STATUS" | jq -e --arg tag "$tag" '.tags | index($tag)' > /dev/null; then
            echo "Tag $tag exists for service $service - OK"
        else
            echo "Tag $tag does not exist for service $service - ERROR"
            exit 1
        fi
    done
done
