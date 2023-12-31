/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

var statusList = [5]string{"critical", "error", "warning", "info", "ok"}

// ServiceStateCollector implements the prometheus.Collector interface
type ServiceStateCollector struct {
	serviceStatusGauge *prometheus.GaugeVec
}

type ServiceStatus struct {
	Timestamp int64
	Status    string
}

func status2int(status string) float64 {
	return float64(indexOf(statusList[:], status))
}

func int2status(status float64) string {
	if status == -1 {
		return "unknown"
	}

	return statusList[int(status)]
}

// NewServiceStateCollector creates a new ServiceStateCollector
func NewServiceStateCollector() *ServiceStateCollector {
	return &ServiceStateCollector{
		serviceStatusGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dashgoat_service_status",
			Help: "Current status of services",
		}, []string{"service"}),
	}
}

// Describe sends the descriptors of each metric over to the provided channel
func (collector *ServiceStateCollector) Describe(ch chan<- *prometheus.Desc) {
	collector.serviceStatusGauge.Describe(ch)
}

// Collect fetches the current state of all services and sends the metrics over to the provided channel
func (collector *ServiceStateCollector) Collect(ch chan<- prometheus.Metric) {
	// Lock the serviceStateList for safe concurrent access
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	// Set the gauge values
	for serviceName, state := range ss.serviceStateList {
		statusValue := status2int(state.Status)
		collector.serviceStatusGauge.WithLabelValues(serviceName).Set(statusValue)
	}
	collector.serviceStatusGauge.Collect(ch)
}

// Delete a service's metric
func deleteServiceMetric(serviceName string) {
	serviceStateCollector.serviceStatusGauge.Delete(prometheus.Labels{"service": serviceName})
}

func queryPrometheus(hours int, serviceID string) ([]ServiceStatus, error) {

	client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090", // Replace with your Prometheus server address
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Prometheus client: %v", err)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(`dashgoat_service_status{service="%s"}`, serviceID)
	result, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
		Start: time.Now().Add(time.Duration(-hours) * time.Hour),
		End:   time.Now(),
		Step:  time.Minute, // Adjust the step to suit your needs
	})
	if err != nil {
		return nil, fmt.Errorf("error querying Prometheus: %v", err)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	matrixVal, ok := result.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("returned result is not a matrix type")
	}

	var statuses []ServiceStatus
	for _, stream := range matrixVal {
		for _, value := range stream.Values {
			statuses = append(statuses, ServiceStatus{
				Timestamp: value.Timestamp.Unix(),
				Status:    int2status(float64(value.Value)),
			})
		}
	}

	return statuses, nil
}

func getMetricsHistory(c echo.Context) error {

	//hours int, serviceID string)
	hours := str2int(c.Param("hours"))
	serviceID := c.Param("serviceid")

	statuses, err := queryPrometheus(hours, serviceID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusOK, statuses)

}
