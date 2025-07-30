package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests.",
}, []string{"method", "path", "status"})

var RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "HTTP request latencies in seconds.",
	Buckets: prometheus.DefBuckets,
}, []string{"method", "path", "status"})

var ClaimSubmissionsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "claim_submissions_total",
	Help: "Total number of claim submissions.",
})

var ClaimReversalsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "claim_reversals_total",
	Help: "Total number of claim reversals.",
})
