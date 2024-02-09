# Dashgoat Application Deployment

This repository contains the Kubernetes configuration files for deploying the dashGoat application.

## Prerequisites

- Kubernetes cluster
- kubectl installed and configured
- Prometheus Operator (for metrics collection)

## Directory Structure

```
.
├── deployment.yaml
├── headless-service.yaml
├── ingress.yaml
├── service-monitor.yaml
├── service.yaml
└── README.md
```

## Deployment Steps

1. Create the namespace:
   ```bash
   kubectl create namespace dashgoat
   ```

2. Apply the Kubernetes configuration files:
   ```bash
   kubectl apply -f deployment.yaml -n dashgoat
   kubectl apply -f service.yaml -n dashgoat
   kubectl apply -f headless-service.yaml -n dashgoat
   ```

The headless service is for discovering its buddies in the same namespace.

3. Update and apply this depending on your ingess :
   ```bash
   kubectl apply -f ingress.yaml -n dashgoat
   ```

4. Apply this for your Prometheus metrics:
   ```bash
   kubectl apply -f service-monitor.yaml -n dashgoat
   ```

## Checking the Status

You can check the status of your resources with the following commands:

- Pods: `kubectl get pods -n dashgoat`
- Services: `kubectl get services -n dashgoat`
- Ingress: `kubectl get ingress -n dashgoat`

## Metrics

Metrics are collected via the `/metrics` endpoint and are scraped by Prometheus using the `ServiceMonitor` configuration.

## Troubleshooting

If you encounter any issues, refer to the [Kubernetes documentation](https://kubernetes.io/docs/home/) or the [Prometheus Operator documentation](https://github.com/prometheus-operator/prometheus-operator).
