/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import "github.com/prometheus/client_golang/prometheus"

var statusList = [5]string{"critical", "error", "warning", "info", "ok"}

// ServiceStateCollector implements the prometheus.Collector interface
type ServiceStateCollector struct {
	serviceStatusGauge *prometheus.GaugeVec
}

func status2int(status string) float64 {
	return float64(indexOf(statusList[:], status))
}

func int2status(status float64) string {
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
