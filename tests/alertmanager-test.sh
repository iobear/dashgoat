#!/bin/bash

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
  echo "Wrong API response for .status"
  exit 1
fi

echo "Test OK, Service: $SERVICE, Host: $HOST, Status: $STATUS"
