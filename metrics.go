package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func metrics() func(w ResponseWriter, r *http.Request) {

	// duration := prometheus.NewHistogram(
	// 	prometheus.HistogramOpts{
	// 		Name:    "scws_request_duration_seconds",
	// 		Help:    "A histogram of latencies for requests.",
	// 		Buckets: []float64{0.1, .25, .5, 1, 2.5, 5, 10},
	// 	},
	// )

	responseSize := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "scws_response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500, 5000, 10000},
		},
	)
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "scws_requests_total",
		Help: "counter of requests",
	})
	prometheus.MustRegister(counter, responseSize)
	return func(w ResponseWriter, r *http.Request) {
		counter.Inc()
		responseSize.Observe(float64(w.Size()))
	}
}
