output "otlp_endpoint" {
  description = "OTLP HTTP endpoint for sending telemetry"
  value       = "http://${var.name}-collector-opentelemetry-collector.${var.namespace}.svc.cluster.local:4318"
}

output "loki_endpoint" {
  description = "Loki endpoint for log queries"
  value       = "http://${var.name}-loki.${var.namespace}.svc.cluster.local:3100"
}

output "grafana_url" {
  description = "Grafana dashboard URL"
  value       = "http://${var.name}-grafana.${var.namespace}.svc.cluster.local:80"
}
