package metrics

import "github.com/zeromicro/go-zero/core/metric"

var RedirectRequestsTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "redirect",
	Name:      "requests_total",
	Help:      "Total short-link resolution requests by outcome.",
	Labels:    []string{"result"},
})

var AnalyticsEventsTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zerolink",
	Subsystem: "analytics",
	Name:      "events_total",
	Help:      "Total visit recording events by outcome.",
	Labels:    []string{"result"},
})
