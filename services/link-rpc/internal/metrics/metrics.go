// Package metrics defines Prometheus business metric counters for link-rpc.
package metrics

import "github.com/zeromicro/go-zero/core/metric"

// RedirectRequestsTotal counts short-link resolution requests by outcome label (hit/miss/disabled/expired/error).
var RedirectRequestsTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "redirect",
	Name:      "requests_total",
	Help:      "Total short-link resolution requests by outcome.",
	Labels:    []string{"result"},
})

// AnalyticsEventsTotal counts visit recording events by outcome label (success/error).
var AnalyticsEventsTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "analytics",
	Name:      "events_total",
	Help:      "Total visit recording events by outcome.",
	Labels:    []string{"result"},
})

// FilterRequestsTotal counts cuckoo filter lookups by outcome label (hit/miss).
var FilterRequestsTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "filter",
	Name:      "requests_total",
	Help:      "Total cuckoo filter lookups by outcome.",
	Labels:    []string{"result"},
})

// CleanupDeletedRowsTotal counts rows deleted by the retention cleanup runner by table label.
var CleanupDeletedRowsTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "cleanup",
	Name:      "deleted_rows_total",
	Help:      "Total rows deleted by the retention cleanup runner.",
	Labels:    []string{"table"},
})
