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


type bmprovisioner_requests_total struct {
	discover prometheus.Counter
	request  prometheus.Counter
	relay_ip *prometheus.CounterVec
}

type bmprovisioner_responses_total struct {
	ack prometheus.Counter
	nack  prometheus.Counter
	offer prometheus.Counter
	relay_ip *prometheus.CounterVec
}

func RequstMetric(reg prometheus.Registerer) *bmprovisioner_requests_total {
	m := &bmprovisioner_requests_total{
		discover: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "req",
			Name: "discover_total",
		}),
		request: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "req",
			Name: "request_total",
		}),
		relay_ip: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relay_req",
			Name: "summ_count",
		},
		[]string{"relay_ip"},
	),
	}
	reg.MustRegister(m.discover, m.request, m.relay_ip)
	return m
}

func ResponseMetric(reg prometheus.Registerer) *bmprovisioner_responses_total {
	m := &bmprovisioner_responses_total{
		ack: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "resp",
			Name: "ack_total",
		}),
		nack: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "resp",
			Name: "nack_total",
		}),
		offer: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "resp",
			Name: "offer_total",
		}),
		relay_ip: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relay_resp",
			Name: "summ_count",
		},
		[]string{"relay_ip"},
	),
	}
	reg.MustRegister(m.ack, m.nack, m.offer, m.relay_ip)
	return m
}

func emulatorRecordMetrics(reg *prometheus.Registry) {
	req := RequstMetric(reg)
	resp := ResponseMetric(reg)
	go func() {
		for {
			// requset
			req.discover.Add(10)
			req.request.Add(10)
			req.relay_ip.With(prometheus.Labels{"relay_ip":"1.1.1.1"}).Add(10)
			req.relay_ip.With(prometheus.Labels{"relay_ip":"2.2.2.2"}).Inc()
			// response
			resp.ack.Inc()
			resp.nack.Inc()
			resp.offer.Inc()
			resp.relay_ip.With(prometheus.Labels{"relay_ip":"1.1.1.1"}).Inc()
			resp.relay_ip.With(prometheus.Labels{"relay_ip":"2.2.2.2"}).Add(10)

			time.Sleep(time.Second)
		}
	}()
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	emulatorRecordMetrics(reg)
	http.Handle("/test", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}
