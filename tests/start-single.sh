#!/bin/bash

echo
echo "-- start-single --"

BASE_URL="http://localhost:2000"

./dashgoat $1 &

PID=$!  # Get the PID
echo $PID > PID-SINGLE-$PID
sleep 1

#Check if PID exists.
if ! kill -0 $PID 2>/dev/null
then
    echo "Error: dashgoat failed to start for instance $INSTANCE"
    exit 1
fi

if [ "$2" = "no-health-test" ]
then
    echo "No /health test"
    exit 0
fi

for i in {1..10}; do
    sleep 1

    STATUS=$(curl -s "$BASE_URL/health")

    if [ "$(echo "$STATUS" | jq -r '.Ready')" = "true" ]; then
        echo "dashGoat ready - OK"
        exit 0
    else
        echo "dashGoat not ready $i"
    fi
done

echo "dashGoat not ready, giving up"
exit 1