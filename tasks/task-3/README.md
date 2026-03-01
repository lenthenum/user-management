# Task 3: Monitoring, Logging & Observability

## 1. Observability Stack Design

* **Structured Logging**: Go backend already uses `slog` (JSON format), which allows Loki to parse fields like `trace_id` and `severity` automatically for better correlation.
* **Logging (Loki)**: Chosen for its native integration with Grafana. Logs are shipped via Promtail.
* **Metrics (Prometheus)**: Scrapes custom `/metrics` endpoints from the Go backend (using Prometheus client library) and standard Kube-state-metrics for cluster health. [WIP]
* **Tracing (Jaeger)**: Provides deep visibility into request latency across microservice boundaries. The Go backend is configured to send spans to the Jaeger collector via OpenTelemetry SDK. [WIP PR#14]
* **Trace Correlation**: By passing the `X-Trace-ID`, logs and traces are linked, allowing you to see the logs for a specific slow request inside the Jaeger UI.

## 2. SLIs and SLOs

* **Availability (SLI)**: Ratio of successful requests (non-5xx) to total requests.
* **Target (SLO)**: 99.9% availability over a 30-day window.

* **Latency (SLI)**: Time taken to process `/api/go/users` requests.
* **Target (SLO)**: 95% of requests completed under 200ms.


## 3. Alerting & Incident Response

* **Critical Alerts**: Should be triggered for Pod Crashes (CrashLoopBackOff), High CPU/Memory usage, and high API error rates.
* **Notification**: Can be configured via Prometheus Alertmanager to send pings to a dedicated Slack channel.
* **Runbooks**: To be included in the `docs/runbooks.md` to provide immediate troubleshooting steps for the "HighErrorRate" alert.

### Alert Example

`prometheus-rules.yaml`:

```yaml
groups:
- name: BackendAlerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "High error rate on backend"
      description: "Backend is returning >5% errors for the last 5 minutes."