package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// bmprovisioner_requests_total{type="Discover, Request", relay_ip="10.0.23.1/24"} 1234
// bmprovisioner_responses_total{type="ACK,Nack,Offer", relay_ip="10.0.23.1/24"} 1234
// bmprovisioner_request_took_seconds{type="Discover, Request"} - prometheus summary -> время обработки одного пакета

// start := time.Now()
// defer func() {
//  summary.WithLabelValues(r.HType()).Observe(time.Since(start).Seconds())
// }()

// type metrics struct {
//     devices prometheus.Gauge
// }

type bmprovisioner_requests_total struct {
	discover prometheus.Gauge
	request  prometheus.Gauge
	// relay_ip *prometheus.CounterVec
}

func RequstMetric(reg prometheus.Registerer) *bmprovisioner_requests_total {
	m := &bmprovisioner_requests_total{
		discover: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "discover_count",
			Help: "123",
		}),
		request: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "request_count",
			Help: "123",
		}),
	}
	reg.MustRegister(m.discover)
	reg.MustRegister(m.request)
	// reg.MustRegister(m.relay_ip)
	return m
}

func recordMetrics(reg *prometheus.Registry) {
	m := RequstMetric(reg)
	go func() {
		for {
			m.discover.Add(10)
			m.request.Add(20)
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	recordMetrics(reg)
	// m := RequstMetric(reg)
	// Set values for the new created metrics.
	// m.discover.Set(10)
	// m.request.Set(20)
	// m.relay_ip.With(prometheus.Labels{"ip":"1.1.1.1"}).Inc()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/test", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}
