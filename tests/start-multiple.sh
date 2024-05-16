#!/bin/bash

echo
echo "-- start multiple --"
echo

MAX=6

# Check if an argument is provided
if [ $# -eq 0 ]
then
    echo "No instance amount argument provided"
    exit 1
fi

INSTANCES=$1

# Validate if the argument is an integer between 1 and $MAX
if [[ $INSTANCES =~ ^[0-9]+$ ]] && [ $INSTANCES -ge 1 ] && [ $INSTANCES -le $MAX ]
then
   echo "Starting $INSTANCES instances"
else
   echo "Error: Argument must be an integer between 1 and $MAX"
   exit 1
fi

start_dashgoat() {
    local INSTANCE=$1
    echo "Starting Instance $INSTANCE"

    if [ -z "$GITHUB_WORKSPACE" ]
    then
        ./dashgoat -ipport ":200$INSTANCE" -dashname "test$INSTANCE" $2 &
    else
        $GITHUB_WORKSPACE/dashgoatbuild/dashgoat -ipport ":200$INSTANCE" -dashname "test$INSTANCE" $2 &
    fi

    PID=$!  # Get the PID
    echo $PID > PID$INSTANCE
    sleep 1

    #Check if PID exists.
    if ! kill -0 $PID 2>/dev/null
    then
        echo "Error: dashgoat failed to start for instance $INSTANCE"
        exit 1
    fi
}

for INSTANCE in $(seq 1 $INSTANCES)
do
    start_dashgoat $INSTANCE
done

##Check if app is ready
for INSTANCE in $(seq 1 $INSTANCES)
do
    READY="false"
    for i in {1..10}
    do
        sleep 1
        STATUS=$(curl -s "http://localhost:200$INSTANCE/health")

        if [ "$(echo "$STATUS" | jq -r '.Ready')" = "true" ]
        then
            echo "dashGoat$INSTANCE ready - OK"
            READY="true"
            break
        else
            echo "dashGoat$INSTANCE not ready $i"
        fi
    done
    if [ $READY = "false" ]
    then
        echo "dashGoat not ready, giving up"
        exit 1
    fi
done