output "otlp_endpoint" {
  description = "OTLP HTTP endpoint (via ADOT Lambda layer)"
  value       = "http://localhost:4318"
}

output "adot_layer_arn" {
  description = "ADOT Lambda layer ARN for adding to functions"
  value       = local.adot_layer_arn
}

output "dashboard_url" {
  description = "CloudWatch dashboard URL"
  value       = "https://${data.aws_region.current.name}.console.aws.amazon.com/cloudwatch/home?region=${data.aws_region.current.name}#dashboards"
}
