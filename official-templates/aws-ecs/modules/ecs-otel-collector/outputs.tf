output "otlp_endpoint" {
  description = "OTLP HTTP endpoint for sending telemetry"
  value       = "http://collector.${var.name}.otel.local:4318"
}

output "collector_endpoint" {
  description = "OTel collector gRPC endpoint"
  value       = "collector.${var.name}.otel.local:4317"
}

output "service_id" {
  description = "ECS service ID"
  value       = aws_ecs_service.this.id
}

output "dashboard_url" {
  description = "CloudWatch dashboard URL"
  value       = "https://${data.aws_region.current.name}.console.aws.amazon.com/cloudwatch/home?region=${data.aws_region.current.name}#dashboards"
}
