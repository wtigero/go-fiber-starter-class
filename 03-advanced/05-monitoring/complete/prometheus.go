package main

import (
	"strconv"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ============ Prometheus Metrics ============

var (
	// Counter: จำนวน requests ทั้งหมด
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Histogram: latency ของ requests
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Gauge: จำนวน requests ที่กำลังประมวลผล
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// Counter: จำนวน errors
	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "error_type"},
	)

	// Gauge: Memory usage
	memoryUsageBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)

	// Counter: Business metrics
	todosCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "todos_created_total",
			Help: "Total number of todos created",
		},
	)

	todosCompleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "todos_completed_total",
			Help: "Total number of todos completed",
		},
	)
)

// ============ Prometheus Middleware ============

// PrometheusMiddleware collects HTTP metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()

		// Track in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		// Track errors
		if c.Response().StatusCode() >= 400 {
			errorType := "client_error"
			if c.Response().StatusCode() >= 500 {
				errorType = "server_error"
			}
			httpErrorsTotal.WithLabelValues(method, path, errorType).Inc()
		}

		return err
	}
}

// ============ Prometheus Handler ============

// SetupPrometheus adds /metrics endpoint
func SetupPrometheus(app *fiber.App) {
	// Expose metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

// ============ Custom Metrics Functions ============

// RecordTodoCreated increments the todos created counter
func RecordTodoCreated() {
	todosCreated.Inc()
}

// RecordTodoCompleted increments the todos completed counter
func RecordTodoCompleted() {
	todosCompleted.Inc()
}

// UpdateMemoryUsage updates the memory usage gauge
func UpdateMemoryUsage(bytes uint64) {
	memoryUsageBytes.Set(float64(bytes))
}
