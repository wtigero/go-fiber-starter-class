# Prometheus Integration

## การติดตั้ง Prometheus

### Docker Compose

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  grafana_data:
```

### prometheus.yml

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'fiber-app'
    static_configs:
      - targets: ['host.docker.internal:3000']
    metrics_path: '/metrics'
```

## Metrics ที่เก็บ

### HTTP Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests |
| `http_request_duration_seconds` | Histogram | Request latency |
| `http_requests_in_flight` | Gauge | Current active requests |
| `http_errors_total` | Counter | Total errors |

### Business Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `todos_created_total` | Counter | Todos created |
| `todos_completed_total` | Counter | Todos completed |
| `app_memory_usage_bytes` | Gauge | Memory usage |

## Grafana Dashboards

### Request Rate
```promql
rate(http_requests_total[5m])
```

### Error Rate
```promql
rate(http_errors_total[5m]) / rate(http_requests_total[5m]) * 100
```

### P99 Latency
```promql
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### Memory Usage
```promql
app_memory_usage_bytes / 1024 / 1024
```

## Usage in Code

```go
package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Add Prometheus middleware
    app.Use(PrometheusMiddleware())

    // Setup /metrics endpoint
    SetupPrometheus(app)

    // Your routes...
    app.Post("/todos", func(c *fiber.Ctx) error {
        // Create todo...
        RecordTodoCreated() // Record metric
        return c.JSON(fiber.Map{"success": true})
    })

    app.Listen(":3000")
}
```

## Alerts (alertmanager)

```yaml
groups:
  - name: fiber-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_errors_total[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P99 latency > 1 second"
```
