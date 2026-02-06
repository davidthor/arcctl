output "otlp_endpoint" {
  description = "OTLP HTTP endpoint (local collector on each EC2 instance)"
  value       = "http://localhost:4318"
}

output "grpc_endpoint" {
  description = "OTLP gRPC endpoint"
  value       = "localhost:4317"
}

output "dashboard_url" {
  description = "CloudWatch dashboard URL"
  value       = "https://${data.aws_region.current.name}.console.aws.amazon.com/cloudwatch/home?region=${data.aws_region.current.name}#dashboards"
}

output "config_parameter" {
  description = "SSM parameter name for OTel config"
  value       = aws_ssm_parameter.otel_config.name
}
