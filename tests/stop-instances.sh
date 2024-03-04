#!/bin/bash

echo
echo "-- stop instances --"
echo

for pid_file in PID*; do
    if [ -f "$pid_file" ]; then
        kill $(cat "$pid_file")
        rm $pid_file
    fi
done
