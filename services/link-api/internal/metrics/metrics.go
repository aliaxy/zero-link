// Package metrics defines Prometheus business metric counters for link-api.
package metrics

import "github.com/zeromicro/go-zero/core/metric"

// AnalyticsDroppedTotal counts analytics events dropped because the worker queue was full.
var AnalyticsDroppedTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "analytics",
	Name:      "dropped_total",
	Help:      "Total analytics events dropped due to a full worker queue.",
	Labels:    []string{},
})

// AnalyticsQueueDepth tracks the instantaneous number of jobs waiting in the analytics worker queue.
var AnalyticsQueueDepth = metric.NewGaugeVec(&metric.GaugeVecOpts{
	Namespace: "zerolink",
	Subsystem: "analytics",
	Name:      "queue_depth",
	Help:      "Current number of jobs waiting in the analytics worker queue.",
	Labels:    []string{},
})
