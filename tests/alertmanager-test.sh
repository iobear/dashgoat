#!/bin/bash

# Alertmanager webhook URL
url="http://localhost:2000/alertmanager/1TvdcoH5RTTTTKLS6CF"

# JSON payload
data='{
  "version": "4",
  "groupKey": "{}:{}",
  "status": "firing",
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

# Send POST request
curl -v -X POST -H "Content-Type: application/json" -d "${data}" ${url}
