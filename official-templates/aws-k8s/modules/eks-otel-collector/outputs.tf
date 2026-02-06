output "otlp_endpoint" {
  description = "OTLP HTTP endpoint"
  value       = "http://${local.name}-opentelemetry-collector.${var.namespace}.svc.cluster.local:4318"
}

output "grpc_endpoint" {
  description = "OTLP gRPC endpoint"
  value       = "${local.name}-opentelemetry-collector.${var.namespace}.svc.cluster.local:4317"
}

output "dashboard_url" {
  description = "CloudWatch dashboard URL"
  value       = "https://${data.aws_region.current.name}.console.aws.amazon.com/cloudwatch/home?region=${data.aws_region.current.name}#dashboards"
}
