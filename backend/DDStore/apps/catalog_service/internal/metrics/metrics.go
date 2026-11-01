// Package metrics records Prometheus metrics for the HTTP API and serves them on a separate
// port, so Prometheus can scrape them on the compose network without /metrics being reachable
// through the service's own routes (the gateway, an ngrok tunnel, ...).
package metrics

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Addr is where /metrics is served. It is not published to the host.
const Addr = ":2112"

var (
	requests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests handled, by route, method and status code.",
	}, []string{"route", "method", "status"})

	duration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time taken to handle an HTTP request, by route and method.",
		Buckets: []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"route", "method"})

	inFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "HTTP requests currently being handled.",
	})
)

// Middleware records each request under its route template (/order/:id, not /order/42),
// so the number of time series stays bounded.
func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		inFlight.Inc()
		defer inFlight.Dec()

		ctx.Next()

		route := ctx.FullPath()
		if route == "" {
			route = "unmatched"
		}
		requests.WithLabelValues(route, ctx.Request.Method, strconv.Itoa(ctx.Writer.Status())).Inc()
		duration.WithLabelValues(route, ctx.Request.Method).Observe(time.Since(start).Seconds())
	}
}

// Serve starts the /metrics endpoint in the background. It also exposes the Go runtime and
// process metrics (memory, GC, goroutines, CPU) that client_golang registers by default.
func Serve() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(Addr, mux); err != nil {
			log.Printf("metrics server stopped: %v", err)
		}
	}()
}
