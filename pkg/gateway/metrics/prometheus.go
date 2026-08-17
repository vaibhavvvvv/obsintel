package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obsintel_requests_total",
			Help: "Total number of requests",
		},
		[]string{"model", "success"},
	)

	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "obsintel_request_latency_ms",
			Help:    "Request latency in milliseconds",
			Buckets: []float64{100, 500, 1000, 2000, 5000, 10000, 15000, 20000, 30000, 50000},
		},
		[]string{"model"},
	)

	CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obsintel_cache_hits_total",
			Help: "Total cache hits vs misses",
		},
		[]string{"result"}, // "hit" or "miss"
	)

	TokensUsed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obsintel_tokens_total",
			Help: "Total tokens consumed",
		},
		[]string{"model", "type"}, // type: "prompt" or "response"
	)
)

func Init() {
	prometheus.MustRegister(RequestTotal)
	prometheus.MustRegister(RequestLatency)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(TokensUsed)
}
