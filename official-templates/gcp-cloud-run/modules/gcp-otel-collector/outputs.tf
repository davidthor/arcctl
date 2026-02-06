output "otlp_endpoint" {
  description = "OTLP gRPC endpoint for the collector"
  value       = "${google_cloud_run_v2_service.otel_collector.uri}:443"
}
