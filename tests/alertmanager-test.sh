#!/bin/bash

echo
echo "-- alertmanager test --"
echo

if [ -z $1 ]; then
  STATUS="firing"
else
  STATUS="resolved"
fi

BASE_URL="http://localhost:2000"
url=$BASE_URL"/alertmanager/1TvdcoH5RTTTTKLS6CF"

data='{
  "version": "4",
  "groupKey": "{}:{}",
  "status": "'$STATUS'",
  "receiver": "webhook",
  "groupLabels": {
    "alertname": "TestAlert"
  },
  "commonLabels": {
    "prometheus": "cluster-0",
    "severity": "error",
    "namespace": "testapp"
  },
  "commonAnnotations": {
    "info": "This is a test alert"
  },
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "TestAlert"
      },
      "annotations": {
        "summary": "This is a test alert"
      }
    }
  ]
}'

echo "Updating dashgoat with Alertmanager alert"
curl -X POST -H "Content-Type: application/json" -d "${data}" ${url}


RESPONSE=$(curl -s "$BASE_URL/status/list")

echo "Validate the response against the expected values"
SERVICE=$(echo "$RESPONSE" | jq '.["cluster-0testapp"].service')

if [[ "$SERVICE" != '"testapp"' ]]; then
  echo "Wrong API response for .service"
  echo "$SERVICE"
  exit 1
fi

HOST=$(echo "$RESPONSE" | jq '.["cluster-0testapp"].host')
if [[ "$HOST" != '"cluster-0"' ]]; then
  echo "Wrong API response for .host"
  echo "$HOST"
  exit 1
fi

STATUS=$(echo "$RESPONSE" | jq '.["cluster-0testapp"].status')
if [[ "$STATUS" != '"error"' ]]; then
  echo "Unexpected API response for .status"
  exit 1
fi

echo "Service: $SERVICE, Host: $HOST, Status: $STATUS"
echo "OK"

echo "Checking probe and change is the same"
PROBE=$(curl -s localhost:2000/status/list | jq '.["cluster-0testapp"].probe')
CHANGE=$(curl -s localhost:2000/status/list | jq '.["cluster-0testapp"].change')

if [[ $PROBE -ne $CHANGE ]]; then

  #1sec difference due to timing
  if [[ $((CHANGE - PROBE)) -eq 1 ]]; then
      echo "OK"
      continue
  fi

  echo "Unexpected API response for .change and .probe"
  echo $PROBE $CHANGE
  exit 1
else
  echo "OK"
fi

sleep 2

echo "Second run, sending Alertmanager alert"
curl -X POST -H "Content-Type: application/json" -d "${data}" ${url}

RESPONSE=$(curl -s "$BASE_URL/status/list")

echo "Validate the response against the expected values"
SERVICE=$(echo "$RESPONSE" | jq '.["cluster-0testapp"].service')

if [[ "$SERVICE" != '"testapp"' ]]; then
  echo "Unexpected API response for .service"
  echo "$SERVICE"
  exit 1
fi
echo "OK"

echo "Checking probe beeing updated, not change"
PROBE=$(curl -s localhost:2000/status/list | jq '.["cluster-0testapp"].probe')
CHANGE=$(curl -s localhost:2000/status/list | jq '.["cluster-0testapp"].change')


if [[ $PROBE -gt $CHANGE ]]; then
    echo "OK"
else
  echo "Unexpected API response for .change"
  echo $PROBE $CHANGE
  exit 1
fi


echo "cleaning up data"
curl -s --request DELETE --url $BASE_URL/service/cluster-0testapp
