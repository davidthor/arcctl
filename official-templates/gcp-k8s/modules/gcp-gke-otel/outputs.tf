output "otlp_endpoint" {
  description = "OTLP gRPC endpoint for the collector"
  value       = "${var.name}-collector.${var.namespace}.svc.cluster.local:4317"
}
