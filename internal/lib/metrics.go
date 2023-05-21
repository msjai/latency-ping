package lib

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CountMinLatencyEndpoint = promauto.NewCounter(prometheus.CounterOpts{
		Name: "count_min_latency_endpoint",
		Help: "counts the number of requests per min latency endponit",
	})

	CountMaxLatencyEndpoint = promauto.NewCounter(prometheus.CounterOpts{
		Name: "count_max_latency_endpoint",
		Help: "counts the number of requests per max latency endponit",
	})

	CountLatencyWebSiteNameEndpoint = promauto.NewCounter(prometheus.CounterOpts{
		Name: "count_latency_website_name_endpoint",
		Help: "counts the number of requests per website name latency endponit",
	})
)
