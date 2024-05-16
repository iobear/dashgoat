#!/bin/bash

echo
echo "-- heartbeat test --"
echo

BASE_URL="http://localhost:2000"
URNKEY="1TvdcoH5RTTTTKLS6CF"

hosts=("web-1" "mail-1" "storage-1")

for host in "${hosts[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/heartbeat/$URNKEY/$host/5/help" \

done

for host in "${hosts[@]}"; do

    echo "Checking probe and change is the same ${host}"
    PROBE=$(curl -s "$BASE_URL/status/${host}heartbeat" | jq '.probe')
    CHANGE=$(curl -s "$BASE_URL/status/${host}heartbeat" | jq '.change')

    if [[ $PROBE -ne $CHANGE ]]; then
        echo $PROBE $CHANGE
        echo "Unexpected API response for .change and .probe ${host}"
        exit 1
    else
        echo "OK"
    fi
done

sleep 1

for host in "${hosts[@]}"; do
    echo "Updating status for service: $service"

    curl -X POST "$BASE_URL/heartbeat/$URNKEY/$host/5/help" \

done


for host in "${hosts[@]}"; do

    echo "Checking probe and change is not the same ${host}"
    PROBE=$(curl -s "$BASE_URL/status/${host}heartbeat" | jq '.probe')
    CHANGE=$(curl -s "$BASE_URL/status/${host}heartbeat" | jq '.change')

    if [[ $PROBE -gt $CHANGE ]]; then
        echo "OK"
    else
        echo $PROBE $CHANGE
        echo "Unexpected API response for .change and .probe ${host}"
        exit 1
    fi
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


echo "cleaning up data"

for host in "${hosts[@]}"; do
    curl -s --request DELETE --url $BASE_URL/service/${host}heartbeat
done